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
<<<<<<< HEAD
	cfg         models.AgentConfig
	Server_Port string
}

func New(server_port string, config models.AgentConfig) *Agent {
=======
	cfg         Config
	Server_Port string
}

func New(server_port string, config Config) *Agent {
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
	return &Agent{
		cfg:         config,
		Server_Port: server_port,
	}
}

func (a *Agent) Work() {
<<<<<<< HEAD
	time.Sleep(time.Millisecond * 2000)

	url := fmt.Sprintf("http://localhost:%s/api/v1/internal/task/", a.Server_Port)
	resp, _ := http.Get(url)

=======
	time.Sleep(time.Millisecond * 200)
	url := fmt.Sprintf("http://localhost:%s/api/v1/internal/task/", a.Server_Port)
	resp, _ := http.Get(url)
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
	type Req struct {
		Task models.Task `json:"task"`
	}
	res := Req{}

	json.NewDecoder(resp.Body).Decode(&res)
<<<<<<< HEAD
	fmt.Println("solving", res.Task)
	guard := make(chan struct{}, a.cfg.ComputingPower)

	for i := 0; i < a.cfg.ComputingPower; i++ {
=======

	guard := make(chan struct{}, a.cfg.Computing_Power)

	for i := 0; i < a.cfg.Computing_Power; i++ {
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
		guard <- struct{}{}
		go func(task *models.Task) {
			a.worker(task)
			<-guard
		}(&res.Task)
	}

	a.Work()
}

func (a *Agent) worker(t *models.Task) {

	time.Sleep(t.OperationTime)
	var result float64
	switch t.Operation {
	case "+":
<<<<<<< HEAD
		timer := time.Duration(a.cfg.TimeAdditionMS) * time.Millisecond
=======
		timer := time.Duration(a.cfg.Addition_Time) * time.Millisecond
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
		time.Sleep(timer)
		t.OperationTime = timer / 1000000
		result = t.Arg1 + t.Arg2
	case "-":
<<<<<<< HEAD
		timer := time.Duration(a.cfg.TimeSubtractionMS) * time.Millisecond
=======
		timer := time.Duration(a.cfg.Substraction_Time) * time.Millisecond
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
		time.Sleep(timer)
		t.OperationTime = timer / 1000000
		result = t.Arg2 - t.Arg1

	case "*":
<<<<<<< HEAD
		timer := time.Duration(a.cfg.TimeMultiplicationMS) * time.Millisecond
=======
		timer := time.Duration(a.cfg.Multiplication_Time) * time.Millisecond
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
		time.Sleep(timer)
		t.OperationTime = timer / 1000000
		result = t.Arg1 * t.Arg2
	case "/":
<<<<<<< HEAD
		timer := time.Duration(a.cfg.TimeDivisionMS) * time.Millisecond
=======
		timer := time.Duration(a.cfg.Division_Time) * time.Millisecond
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
		time.Sleep(timer)
		t.OperationTime = timer / 1000000
		result = t.Arg1 / t.Arg2
	default:
		result = 0
	}

	t.Value = result
	t.Status = "Resolved"
<<<<<<< HEAD
	url := fmt.Sprintf("http://localhost:%s/api/v1/internal/post_task/", a.Server_Port)
=======
	url := fmt.Sprintf("http://localhost:%s/api/v1/internal/task/", a.Server_Port)
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868

	type Request struct {
		Task models.Task `json:"task"`
	}
	req := Request{Task: *t}
	ans_bytes, _ := json.Marshal(req)
	r := bytes.NewReader(ans_bytes)
<<<<<<< HEAD
	//	fmt.Println(url)
=======

>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
	http.Post(url, "application/json", r)
}
