package application

import (
	"encoding/json"
	"fmt"
	calculator "lol/pkg"
	"net/http"
	"os"
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
	Expression string  `json:"Expression"`
	Result     float64 `json:"Result"`
	Error      string  `json:"Error"`
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
			fmt.Fprintf(w, string(ans_bytes))
		}
	}

}

func (a *Application) RunServer() error {
	http.HandleFunc("/", CalculationHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
