package orchestrator

import (
	"bytes"
	"calc/pkg/calculator"
	"calc/pkg/config"
	logger "calc/pkg/logger"
	"calc/server/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	"go.uber.org/zap"
)

var expressions models.ExpressionQueue

var tasks models.TaskQueue

var db *sql.DB

var global_userID int = 1

// type Config struct {
// 	Addr string
// }

type ServerCorrectAnswer struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

type ServerErrorAnswer struct {
	Expression string `json:"expression"`
	Error      string `json:"error"`
}

type Application struct {
	cfg config.Config
}

func loggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if r.URL.Path != "/api/v1/internal/task/" && r.URL.Path != "/api/v1/internal/post_task/" {

				duration := time.Since(start)
				next.ServeHTTP(w, r)
				logger.Info("HTTP request",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),

					zap.Duration("duration", duration),

					zap.String("body", string(bodyBytes)),
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *Application) generateToken(userID int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(app.cfg.Jwt_expiration).Unix(),
		"jti":     uuid.New().String(),
	})
	return token.SignedString([]byte(app.cfg.Jwt_secret_key))
}

func (a *Application) validateToken(authHeader string) (int, error) {

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return 0, fmt.Errorf("authorization header format must be 'Bearer {token}'")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return 0, errors.New("Token contains an invalid number of segments")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.cfg.Jwt_secret_key), nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64) // JWT числа всегда float64
		if !ok {
			return 0, fmt.Errorf("invalid user_id in token")
		}
		return int(userID), nil
	}

	return 0, fmt.Errorf("invalid token claims")
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./calculator.db?_busy_timeout=1000&_journal_mode=WAL")

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	type Expression struct {
		Id            string    `json:"id"`
		Status        string    `json:"status"`
		Result        float64   `json:"result"`
		Value         string    `json:"value"`
		PostfixString []string  `json:"-"`
		WaitforSolve  bool      `json:"-"`
		Stack         []float64 `json:"-"`
		SavedIndex    int       `json:"-"`
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);


	
	CREATE TABLE IF NOT EXISTS expressions (
		id VARCHAR(36) PRIMARY KEY,
		user_id INTEGER NOT NULL,
		status VARCHAR(20) NOT NULL,
		result FLOAT,
		value TEXT NOT NULL,
		postfix_string TEXT,
		wait_for_solve BOOLEAN DEFAULT FALSE,
		stack TEXT,
		saved_index INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
);`)
	if err != nil {
		log.Fatal(err)
	}
}

func New(cfg config.Config) *Application {
	return &Application{
		cfg: cfg,
	}
}

func (a *Application) PrintAllExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	enableCORS(w, r)
	user_id, err := a.validateToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	global_userID = user_id
	// Получаем выражения из базы данных для текущего пользователя
	rows, err := db.Query(`
        SELECT id, status, result, value, created_at, updated_at 
        FROM expressions 
        WHERE user_id = ? 
        ORDER BY created_at DESC`,
		user_id,
	)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Собираем результаты
	var expressions []struct {
		Id        string    `json:"id"`
		Status    string    `json:"status"`
		Result    float64   `json:"result"`
		Value     string    `json:"value"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	for rows.Next() {
		var expr struct {
			Id        string    `json:"id"`
			Status    string    `json:"status"`
			Result    float64   `json:"result"`
			Value     string    `json:"value"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}

		err := rows.Scan(
			&expr.Id,
			&expr.Status,
			&expr.Result,
			&expr.Value,
			&expr.CreatedAt,
			&expr.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning expressions: "+err.Error(), http.StatusInternalServerError)
			return
		}

		expressions = append(expressions, expr)
	}

	// Проверяем ошибки итерации
	if err = rows.Err(); err != nil {
		http.Error(w, "Error during rows iteration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(expressions); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func PrintAllExpressionsHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	enableCORS(w, r)

	ans_bytes, err := json.Marshal(expressions)
	if err != nil {
		http.Error(w, "error in creating expressions json", http.StatusInternalServerError)
	}
	fmt.Fprintln(w, string(ans_bytes))
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

}

func (a *Application) NewExpressionHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	enableCORS(w, r)
	type Request struct {
		Expression string `json:"expression"`
	}
	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "error in parsing json", http.StatusInternalServerError)
	} else {
		type Answer struct {
			Id string `json:"id"`
		}

		//id, err := expressions.AddExpression(req.Expression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}

		user_id, err := a.validateToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}

		global_userID = user_id
		postfix, err := calculator.CreatePostfix(req.Expression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		expr := &models.Expression{Id: uuid.New().String(), Value: req.Expression, Status: "Waiting", Result: 0.0, PostfixString: postfix, Stack: []float64{}, SavedIndex: 0}
		postfixStr, _ := json.Marshal(expr.PostfixString)
		stackStr, _ := json.Marshal(expr.Stack)
		_, err = db.Exec(`
        INSERT INTO expressions (
            id, user_id, status, result, value, 
            postfix_string, wait_for_solve, stack, saved_index
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			expr.Id,
			user_id,
			expr.Status,
			expr.Result,
			expr.Value,
			string(postfixStr),
			expr.WaitForSolve,
			string(stackStr),
			expr.SavedIndex,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		server_ans := Answer{Id: expr.Id}

		id, err := expressions.AddExpression(req.Expression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		server_ans = Answer{Id: id}

		ans_bytes, _ := json.Marshal(server_ans)
		fmt.Fprintln(w, string(ans_bytes))
	}
}

func (a *Application) TaskCreator() {
	time.Sleep(time.Second * 5)

	rows, err := db.Query(`
        SELECT id, status, result, value, postfix_string, 
               wait_for_solve, stack, saved_index
        FROM expressions 
        WHERE user_id = ? AND status NOT IN ('Solved', 'Error in expression')`,
		global_userID,
	)
	if err != nil {
		log.Printf("Failed to get expressions: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var expr models.Expression
		var postfixStr, stackStr string
		err := rows.Scan(
			&expr.Id,
			&expr.Status,
			&expr.Result,
			&expr.Value,
			&postfixStr,
			&expr.WaitForSolve,
			&stackStr,
			&expr.SavedIndex,
		)
		if err != nil {
			log.Printf("Failed to scan expression: %v", err)
			continue
		}
		json.Unmarshal([]byte(postfixStr), &expr.PostfixString)
		json.Unmarshal([]byte(stackStr), &expr.Stack)
		fmt.Println(expr.Id, expr.SavedIndex, expr.PostfixString, expr.Stack)

		if expr.Status != "Solved" && expr.Status != "Error in expression" && !expr.WaitForSolve {
			if expr.SavedIndex == len(expr.PostfixString) && len(expr.Stack) != 0 {
				expr.Result = expr.Stack[0]
				db.Exec(`
                UPDATE expressions 
                SET status = 'Solved', 
                    result = ?,
                    updated_at = CURRENT_TIMESTAMP
                WHERE id = ?`,
					expr.Result,
					expr.Id,
				)
			} else {
				for i := expr.SavedIndex; i < len(expr.PostfixString); i++ {
					if !expr.WaitForSolve {
						val := expr.PostfixString[i]
						conv_val, err := strconv.ParseFloat(val, 64)
						if err == nil {
							expr.Stack = append(expr.Stack, conv_val)
							stackStr, _ := json.Marshal(expr.Stack)
							db.Exec(`
                                UPDATE expressions 
                                SET stack = ?,
                                    updated_at = CURRENT_TIMESTAMP
                                WHERE id = ?`,
								string(stackStr),
								expr.Id,
							)
						} else if calculator.IsOperator(val) {
							fmt.Println("NICE NICE NICE")

							fir_pop_item := expr.Stack[len(expr.Stack)-1]
							expr.Stack = expr.Stack[:len(expr.Stack)-1]
							sec_pop_item := expr.Stack[len(expr.Stack)-1]
							expr.Stack = expr.Stack[:len(expr.Stack)-1]

							if fir_pop_item == 0 && val == "/" {
								db.Exec(`
                                    UPDATE expressions 
                                    SET status = 'Error in expression',
                                        updated_at = CURRENT_TIMESTAMP
                                    WHERE id = ?`,
									expr.Id,
								)
							}

							tasks.NewTask(expr.Id, fir_pop_item, sec_pop_item, val)
							expr.WaitForSolve = true
							expr.SavedIndex = i + 1
							stackStr, _ := json.Marshal(expr.Stack)
							db.Exec(`
                                UPDATE expressions 
                                SET stack = ?,
                                    saved_index = ?,
                                    wait_for_solve = true,
                                    updated_at = CURRENT_TIMESTAMP
                                WHERE id = ?`,
								string(stackStr),
								expr.SavedIndex,
								expr.Id,
							)
						}
					}
				}
			}
		}
	}
	a.TaskCreator()
}

func PrintAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	enableCORS(w, r)

	ans_bytes, _ := json.Marshal(tasks)
	fmt.Fprintln(w, string(ans_bytes))

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {

		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		"INSERT INTO users (email, username, password) VALUES (?, ?, ?)",
		creds.Email, creds.Username, creds.Password,
	)

	if err != nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	r.Body.Close() // Явно закрываем тело

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.Unmarshal(bodyBytes, &creds); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	var userID int
	var dbPassword string
	err = db.QueryRow(
		"SELECT id, password FROM users WHERE email = ?",
		creds.Email,
	).Scan(&userID, &dbPassword)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if dbPassword != creds.Password {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	token, err := app.generateToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"token":   token,
		"user_id": userID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}

func PrintExpressionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	enableCORS(w, r)
	var creds struct {
		Expression_id string `json:"expression_id"`
		User_id       string `json:"user_id"`
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	r.Body.Close() // Явно закрываем тело
	if err := json.Unmarshal(bodyBytes, &creds); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}
	fmt.Println("CREDS", creds.Expression_id, creds.User_id)
	var expr models.Expression
	err = db.QueryRow(
		"SELECT id, status, result, value FROM expressions WHERE id = ? AND user_id = ?",
		creds.Expression_id,
		creds.User_id,
	).Scan(
		&expr.Id,
		&expr.Status,
		&expr.Result,
		&expr.Value,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	ans_bytes, _ := json.Marshal(expr)
	fmt.Fprintln(w, string(ans_bytes))
}

func TaskSendHandler(w http.ResponseWriter, r *http.Request) {

	for _, task := range tasks.Tasks {
		if task.Status == "Unresolved" && !task.IsSolving {
			task.IsSolving = true
			type Response struct {
				Task models.Task `json:"task"`
			}

			resp := Response{Task: *task}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
}

func TaskSolveHandler(w http.ResponseWriter, r *http.Request) {
	type taskReq struct {
		Task *models.Task `json:"task"`
	}

	req := taskReq{}
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "error in parsing json", http.StatusUnprocessableEntity)
	}

	for _, task := range tasks.Tasks {
		if task.Id == req.Task.Id {
			task.OperationTime = req.Task.OperationTime
			task.Status = "Resolved"
			task.Value = req.Task.Value
		}
	}

	rows, err := db.Query(`
        SELECT id, status, result, value, postfix_string, 
               wait_for_solve, stack, saved_index
        FROM expressions 
        WHERE user_id = ? AND status NOT IN ('Solved', 'Error in expression')`,
		global_userID,
	)
	if err != nil {
		log.Printf("Failed to get expressions: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {

		var expr models.Expression
		var postfixStr, stackStr string
		err := rows.Scan(
			&expr.Id,
			&expr.Status,
			&expr.Result,
			&expr.Value,
			&postfixStr,
			&expr.WaitForSolve,
			&stackStr,
			&expr.SavedIndex,
		)
		if err != nil {
			log.Printf("Failed to scan expression: %v", err)
			continue
		}
		json.Unmarshal([]byte(postfixStr), &expr.PostfixString)
		json.Unmarshal([]byte(stackStr), &expr.Stack)
		if expr.Id == req.Task.ExpressionId && expr.WaitForSolve {

			expr.WaitForSolve = false
			expr.Stack = append(expr.Stack, req.Task.Value)
			stackStr, _ := json.Marshal(expr.Stack)
			db.Exec(`
					UPDATE expressions 
					SET stack = ?,
					wait_for_solve = ?,
					updated_at = CURRENT_TIMESTAMP
					WHERE id = ?`,
				string(stackStr),
				expr.WaitForSolve,
				expr.Id,
			)
		}
	}
	for _, expr := range expressions.Expressions {
		if expr.Id == req.Task.ExpressionId && expr.WaitForSolve {
			expr.WaitForSolve = false
			expr.Stack = append(expr.Stack, req.Task.Value)

		}
	}

}

func (a *Application) RunServer() error {
	initDB()
	defer db.Close()
	r := mux.NewRouter()
	logger := logger.SetupLogger()
	staticDir := "./static" // Путь к папке со статикой
	staticHandler := http.FileServer(http.Dir(staticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем правильные MIME-типы
		switch {
		case strings.HasSuffix(r.URL.Path, ".css"):
			w.Header().Set("Content-Type", "text/css")
		case strings.HasSuffix(r.URL.Path, ".js"):
			w.Header().Set("Content-Type", "application/javascript")
		case strings.HasSuffix(r.URL.Path, ".ico"):
			w.Header().Set("Content-Type", "image/x-icon")
		}

		staticHandler.ServeHTTP(w, r)
	})))

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/login/", a.LoginHandler).Methods("POST")
	api.HandleFunc("/register/", RegisterHandler).Methods("POST")
	api.HandleFunc("/calculate/", a.NewExpressionHandler).Methods("POST")
	api.HandleFunc("/expressions/", a.PrintAllExpressionsHandler).Methods("GET")
	api.HandleFunc("/get_expression/", PrintExpressionHandler).Methods("POST")
	api.HandleFunc("/tasks/", PrintAllTasksHandler).Methods("POST")

	api.HandleFunc("/internal/task/", TaskSendHandler).Methods("GET")        // Агент получает задачу
	api.HandleFunc("/internal/post_task/", TaskSolveHandler).Methods("POST") // Агент возвращает резул

	// 3. Все остальные GET запросы -> index.html
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Только GET запросы и не начинающиеся с /api или /static
		if r.Method == "GET" &&
			!strings.HasPrefix(r.URL.Path, "/api") &&
			!strings.HasPrefix(r.URL.Path, "/static") {

			http.ServeFile(w, r, "./static/index.html")
		} else {
			fmt.Println("Not found:", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	})
	logger.Info("HTTP request",
		zap.String("server status", "started"),
	)
	r.Use(loggingMiddleware(logger))
	return http.ListenAndServe(":"+a.cfg.Server_Port, r)

}
