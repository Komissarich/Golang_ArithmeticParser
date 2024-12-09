package application

import (
	calculator "calc/pkg"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

func setupLogger() *zap.Logger {

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Ошибка настройки логгера: %v\n", err)
	}

	return logger
}

func loggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			duration := time.Since(start)
			next.ServeHTTP(w, r)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)

		})
	}
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type ServerAnswer struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
	Error      string  `json:"error"`
}

func CalculationHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		server_ans := ServerAnswer{Expression: req.Expression, Result: 0, Error: "error in parsing json"}
		ans_bytes, _ := json.Marshal(server_ans)
		http.Error(w, string(ans_bytes), http.StatusBadRequest)

	} else {
		result, err := calculator.Calc(req.Expression)

		var error_text = ""

		if err != nil {
			error_text = err.Error()
		}

		server_ans := ServerAnswer{Expression: req.Expression, Result: result, Error: error_text}
		ans_bytes, _ := json.Marshal(server_ans)

		if err != nil {
			http.Error(w, string(ans_bytes), http.StatusInternalServerError)

		} else {
			fmt.Fprintln(w, string(ans_bytes))
		}
	}

}

func (a *Application) RunServer() error {
	r := mux.NewRouter()

	logger := setupLogger()

	r.Use(loggingMiddleware(logger))
	r.HandleFunc("/", CalculationHandler)

	http.Handle("/", r)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
