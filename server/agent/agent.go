package agent

import (
	"bytes"
	"calc/server/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Addition_Time       int64 `yaml:"TIME_ADDITION_MS" `
	Substraction_Time   int   `TIME_SUBTRACTION_MS" `
	Multiplication_Time int   `yaml:"TIME_MULTIPLICATIONS_MS" `
	Division_Time       int   `yaml:"TIME_DIVISIONS_MS"`
	Computing_Power     int   `yaml:"COMPUTING_POWER"`
}
type Agent struct {
	cfg         models.AgentConfig
	Server_Port string
}

func New(server_port string, config models.AgentConfig) *Agent {

	return &Agent{
		cfg:         config,
		Server_Port: server_port,
	}
}

func (a *Agent) Work() {

	time.Sleep(time.Millisecond * 2000)

	url := fmt.Sprintf("http://localhost:%s/api/v1/internal/task/", a.Server_Port)
	resp, _ := http.Get(url)

	time.Sleep(time.Millisecond * 200)

	type Req struct {
		Task models.Task `json:"task"`
	}
	res := Req{}

	json.NewDecoder(resp.Body).Decode(&res)

	for i := 0; i < a.cfg.ComputingPower; i++ {

		guard := make(chan struct{}, a.cfg.ComputingPower)

		for i := 0; i < a.cfg.ComputingPower; i++ {

			guard <- struct{}{}
			go func(task *models.Task) {
				a.worker(task)
				<-guard
			}(&res.Task)
		}

		a.Work()
	}
}

func (a *Agent) worker(t *models.Task) {

	time.Sleep(t.OperationTime)
	var result float64
	switch t.Operation {
	case "+":

		timer := time.Duration(a.cfg.TimeAdditionMS) * time.Millisecond

		time.Sleep(timer)
		t.OperationTime = timer / 1000000
		result = t.Arg1 + t.Arg2
	case "-":

		timer := time.Duration(a.cfg.TimeSubtractionMS) * time.Millisecond

		time.Sleep(timer)
		t.OperationTime = timer / 1000000
		result = t.Arg2 - t.Arg1

	case "*":

		timer := time.Duration(a.cfg.TimeMultiplicationMS) * time.Millisecond

		time.Sleep(timer)
		t.OperationTime = timer / 1000000
		result = t.Arg1 * t.Arg2
	case "/":

		timer := time.Duration(a.cfg.TimeDivisionMS) * time.Millisecond

		time.Sleep(timer)
		t.OperationTime = timer / 1000000
		result = t.Arg1 / t.Arg2
	default:
		result = 0
	}

	t.Value = result
	t.Status = "Resolved"

	url := fmt.Sprintf("http://localhost:%s/api/v1/internal/post_task/", a.Server_Port)

	type Request struct {
		Task models.Task `json:"task"`
	}
	req := Request{Task: *t}
	ans_bytes, _ := json.Marshal(req)
	r := bytes.NewReader(ans_bytes)

	//	fmt.Println(url)

	http.Post(url, "application/json", r)
}
