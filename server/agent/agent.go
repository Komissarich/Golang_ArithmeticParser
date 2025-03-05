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
	cfg         Config
	Server_Port string
}

func New(server_port string, config Config) *Agent {
	return &Agent{
		cfg:         config,
		Server_Port: server_port,
	}
}

func (a *Agent) Work() {
	time.Sleep(time.Millisecond * 200)
	url := fmt.Sprintf("http://localhost:%s/api/v1/internal/task/", a.Server_Port)
	resp, _ := http.Get(url)
	type Req struct {
		Task models.Task `json:"task"`
	}
	res := Req{}

	json.NewDecoder(resp.Body).Decode(&res)
	//fmt.Println(res.Task)
	//	mutex.Unlock()
	// if err != nil {
	// 	//	fmt.Println("err")
	// } else {
	// 	fmt.Println(res.Task)

	guard := make(chan struct{}, a.cfg.Computing_Power)

	for i := 0; i < a.cfg.Computing_Power; i++ {
		guard <- struct{}{} // would block if guard channel is already filled
		go func(task *models.Task) {
			a.worker(task)
			<-guard
		}(&res.Task)
	}

	a.Work()
}

func (a *Agent) worker(t *models.Task) {
	//fmt.Println("doing work on", i)
	time.Sleep(t.OperationTime)
	var result float64
	switch t.Operation {
	case "+":
		timer := time.Duration(a.cfg.Addition_Time) * time.Millisecond
		time.Sleep(timer)
		t.OperationTime = time.Duration(timer.Seconds())
		result = t.Arg1 + t.Arg2
	case "-":
		timer := time.Duration(a.cfg.Substraction_Time) * time.Millisecond
		time.Sleep(timer)
		t.OperationTime = time.Duration(timer.Seconds())
		result = t.Arg2 - t.Arg1

	case "*":
		timer := time.Duration(a.cfg.Multiplication_Time) * time.Millisecond
		time.Sleep(timer)
		t.OperationTime = time.Duration(timer.Seconds())
		result = t.Arg1 * t.Arg2
	case "/":
		timer := time.Duration(a.cfg.Division_Time) * time.Millisecond
		time.Sleep(timer)
		t.OperationTime = time.Duration(timer.Seconds())
		result = t.Arg1 / t.Arg2
	default:
		result = 0
	}

	t.Value = result
	t.Status = "Resolved"
	url := fmt.Sprintf("http://localhost:%s/api/v1/internal/task/", a.Server_Port)

	type Request struct {
		Task models.Task `json:"task"`
	}
	req := Request{Task: *t}
	ans_bytes, _ := json.Marshal(req)
	r := bytes.NewReader(ans_bytes)

	http.Post(url, "application/json", r)
}
