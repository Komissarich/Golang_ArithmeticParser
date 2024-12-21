package application

import (
	"bytes"
	calculator "calc/pkg"
	"encoding/json"
	"fmt"
	"io"
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

			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			duration := time.Since(start)
			next.ServeHTTP(w, r)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
				zap.String("body", string(bodyBytes)),
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

type ServerCorrectAnswer struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

type ServerErrorAnswer struct {
	Expression string `json:"expression"`
	Error      string `json:"error"`
}

func CalculationHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		server_ans := ServerErrorAnswer{Expression: req.Expression, Error: "error in parsing json"}
		ans_bytes, _ := json.Marshal(server_ans)
		http.Error(w, string(ans_bytes), http.StatusUnprocessableEntity)

	} else {
		result, err := calculator.Calc(req.Expression)

		if err != nil {
			error_text := err.Error()
			server_ans := ServerErrorAnswer{Expression: req.Expression, Error: error_text}
			ans_bytes, _ := json.Marshal(server_ans)
			http.Error(w, string(ans_bytes), http.StatusInternalServerError)

		} else {

			server_ans := ServerCorrectAnswer{Expression: req.Expression, Result: result}
			ans_bytes, _ := json.Marshal(server_ans)
			fmt.Fprintln(w, string(ans_bytes))
		}

	}

}

func (a *Application) RunServer() error {
	r := mux.NewRouter()

	logger := setupLogger()

	r.Use(loggingMiddleware(logger))
	r.HandleFunc("/api/v1/calculate", CalculationHandler)

	http.Handle("/api/v1/calculate", r)
	logger.Info("HTTP request",
		zap.String("server status", "started"),
	)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
