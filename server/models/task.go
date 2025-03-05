package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id            string        `json:"id"`
	ExpressionId  string        `json:"expression_id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	Value         float64       `json:"value"`
	Status        string        `json:"status"`
	OperationTime time.Duration `json:"operation_time"`
	IsSolving     bool          `json:"-"`
}

type TaskQueue struct {
	Tasks []*Task
}

func (t *TaskQueue) NewTask(exprid string, arg1 float64, arg2 float64, op string) {

	task := &Task{Id: uuid.NewString(), IsSolving: false, ExpressionId: exprid, Arg1: arg1, Arg2: arg2, Operation: op, OperationTime: time.Duration(0), Value: 0, Status: "Unresolved"}
	t.Tasks = append(t.Tasks, task)

}
