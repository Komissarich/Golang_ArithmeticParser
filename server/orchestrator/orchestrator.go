package orchestrator

import (
	"bytes"
	"calc/pkg/calculator"
	"calc/pkg/config"
	logger "calc/pkg/logger"
	"calc/server/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var expressions models.ExpressionQueue

var tasks models.TaskQueue

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
			next.ServeHTTP(w, r)
			if r.URL.Path != "/api/v1/internal/task/" {
				duration := time.Since(start)
				next.ServeHTTP(w, r)
				logger.Info("HTTP request",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Duration("duration", duration),
					zap.String("body", string(bodyBytes)),
				)
			}
		})
	}
}

func New(cfg config.Config) *Application {
	return &Application{
		cfg: cfg,
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
	}
}

func NewExpressionHandler(w http.ResponseWriter, r *http.Request) {

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
		id, err := expressions.AddExpression(req.Expression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		server_ans := Answer{Id: id}
		ans_bytes, _ := json.Marshal(server_ans)
		fmt.Fprintln(w, string(ans_bytes))
	}
}

func (a *Application) TaskCreator() {
	time.Sleep(time.Second * 5)

	for _, expr := range expressions.Expressions {
		if expr.Status != "Solved" && expr.Status != "Error in expression" && !expr.WaitforSolve {

			if expr.SavedIndex == len(expr.PostfixString) {
				expr.Status = "Solved"
				expr.Result = expr.Stack[0]
			} else {

				for i := expr.SavedIndex; i < len(expr.PostfixString); i++ {

					if !expr.WaitforSolve {
						val := expr.PostfixString[i]
						conv_val, err := strconv.ParseFloat(val, 64)
						if err == nil {
							expr.Stack = append(expr.Stack, conv_val)
						} else if calculator.IsOperator(val) {

							fir_pop_item := expr.Stack[len(expr.Stack)-1]
							expr.Stack = expr.Stack[:len(expr.Stack)-1]
							sec_pop_item := expr.Stack[len(expr.Stack)-1]
							expr.Stack = expr.Stack[:len(expr.Stack)-1]

							if fir_pop_item == 0 && val == "/" {
								expr.Status = "Error in expression"
							}

							tasks.NewTask(expr.Id, fir_pop_item, sec_pop_item, val)
							expr.WaitforSolve = true

							expr.SavedIndex = i + 1

						}
					}
				}
			}
		}
	}
	a.TaskCreator()
}

func PrintTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	enableCORS(w, r)

	vars := mux.Vars(r)
	id := vars["task_id"]
	for _, task := range tasks.Tasks {
		if task.Id == id {
			ans_bytes, _ := json.Marshal(task)
			fmt.Fprintln(w, string(ans_bytes))
			return
		}
	}
	http.Error(w, "not found", http.StatusNotFound)
}

func PrintAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	enableCORS(w, r)

	ans_bytes, _ := json.Marshal(tasks)
	fmt.Fprintln(w, string(ans_bytes))

}

func PrintExpressionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	enableCORS(w, r)
	vars := mux.Vars(r)
	id := vars["expr_id"]
	for _, expr := range expressions.Expressions {
		if expr.Id == id {
			ans_bytes, _ := json.Marshal(expr)
			fmt.Fprintln(w, string(ans_bytes))
			return
		}
	}
	http.Error(w, "not found", http.StatusNotFound)
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
	for _, expr := range expressions.Expressions {
		if expr.Id == req.Task.ExpressionId && expr.WaitforSolve {
			expr.WaitforSolve = false
			expr.Stack = append(expr.Stack, req.Task.Value)

		}
	}

}

func (a *Application) RunServer() error {
	r := mux.NewRouter()

	logger := logger.SetupLogger()

	r.HandleFunc("/api/v1/calculate/", NewExpressionHandler).Methods("POST")
	r.HandleFunc("/api/v1/expressions/", PrintAllExpressionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{expr_id}", PrintExpressionHandler).Methods("GET")
	r.HandleFunc("/api/v1/tasks/", PrintAllTasksHandler).Methods("GET")
	r.HandleFunc("/api/v1/tasks/{task_id}", PrintTaskHandler).Methods("GET")
	r.HandleFunc("/api/v1/internal/task/", TaskSendHandler).Methods("GET")
	r.HandleFunc("/api/v1/internal/task/", TaskSolveHandler).Methods("POST")

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	logger.Info("HTTP request",
		zap.String("server status", "started"),
	)
	r.Use(loggingMiddleware(logger))
	return http.ListenAndServe(":"+a.cfg.Server_Port, r)
}
