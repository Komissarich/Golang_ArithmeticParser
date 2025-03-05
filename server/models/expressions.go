package models

import (
	"calc/pkg/calculator"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

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

type ExpressionQueue struct {
	Expressions []*Expression `json:"expressions"`
}

func (e *ExpressionQueue) AddExpression(expression string) (string, error) {
	postfix, err := calculator.CreatePostfix(expression)

	expr := &Expression{Id: uuid.New().String(), Value: expression, Status: "Waiting", Result: 0.0, PostfixString: postfix, Stack: []float64{}, SavedIndex: 0}
	if err != nil || len(postfix) == 0 {
		expr.Status = "Error in expression"
		e.Expressions = append(e.Expressions, expr)
		return expr.Id, err
	}
	e.Expressions = append(e.Expressions, expr)
	return expr.Id, nil
}

func (e *ExpressionQueue) ChangeExpressionStatus(id string) {
	for _, expr := range e.Expressions {
		if expr.Id == id {
			expr.Status = "Solved"
		}
	}
}

func (e *ExpressionQueue) Print(w http.ResponseWriter) {
	ans_bytes, _ := json.Marshal(e)

	fmt.Fprintln(w, string(ans_bytes))
}
