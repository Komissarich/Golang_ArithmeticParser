package application_test

import (
	"calc/application"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestHandlerGoodRequestCase(t *testing.T) {

	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult application.ServerCorrectAnswer
	}{
		{
			name:           "simple",
			expression:     `{"expression":"1+1"}`,
			expectedResult: application.ServerCorrectAnswer{Expression: "1+1", Result: 2},
		},
		{
			name:           "priority",
			expression:     `{"expression":"(2+2)*2"}`,
			expectedResult: application.ServerCorrectAnswer{Expression: "(2+2)*2", Result: 8},
		},
		{
			name:           "priority2",
			expression:     `{"expression":"2+2-2"}`,
			expectedResult: application.ServerCorrectAnswer{Expression: "2+2-2", Result: 2},
		},
		{
			name:           "division",
			expression:     `{"expression":"1/2"}`,
			expectedResult: application.ServerCorrectAnswer{Expression: "1/2", Result: 0.5},
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.expression)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", reader)
			w := httptest.NewRecorder()
			application.CalculationHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			var ans application.ServerCorrectAnswer
			err := json.NewDecoder(res.Body).Decode(&ans)
			if err != nil {
				t.Fatalf("successful case %s returns error", testCase.expression)
			} else if res.StatusCode != 200 {
				t.Fatalf("%s test returns wrong status code", testCase.expression)
			} else if ans != testCase.expectedResult {
				t.Fatalf("%+v\n should be equal %+v\n", ans, testCase.expectedResult)
			} else {
				fmt.Println("great test!")
			}
		})
	}

}

func TestRequestHandlerBadRequestCase(t *testing.T) {
	testCasesFail := []struct {
		name           string
		expression     string
		expectedResult application.ServerErrorAnswer
	}{
		{
			name:           "operators",
			expression:     `{"expression":"1+*1"}`,
			expectedResult: application.ServerErrorAnswer{Expression: "1+*1", Error: "internal server error: repeating operators"},
		},
		{
			name:           "operator2",
			expression:     `{"expression":"2+2**2"}`,
			expectedResult: application.ServerErrorAnswer{Expression: "2+2**2", Error: "internal server error: repeating operators"},
		},
		{
			name:           "brackets",
			expression:     `{"expression":"((2+2-*(2"}`,
			expectedResult: application.ServerErrorAnswer{Expression: "((2+2-*(2", Error: "internal server error: invalid brackets"},
		},
		{
			name:           "empty",
			expression:     `{"expression":""}`,
			expectedResult: application.ServerErrorAnswer{Expression: "", Error: "internal server error: empty string"},
		},
		{
			name:           "invalid expression",
			expression:     `{"expression":"a+b"}`,
			expectedResult: application.ServerErrorAnswer{Expression: "a+b", Error: "expression is not valid"},
		},
		{
			name:           "division by zero",
			expression:     `{"expression":"2/0"}`,
			expectedResult: application.ServerErrorAnswer{Expression: "2/0", Error: "internal server error: division by zero"},
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.expression)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", reader)
			w := httptest.NewRecorder()
			application.CalculationHandler(w, req)
			res := w.Result()
			defer res.Body.Close()

			var ans application.ServerErrorAnswer
			json.NewDecoder(res.Body).Decode(&ans)

			if ans != testCase.expectedResult {
				t.Fatalf("%+v\n should be equal %+v\n", ans, testCase.expectedResult)
			} else if res.StatusCode != 500 {
				t.Fatalf("%s test returns wrong status code", testCase.expression)
			} else {
				fmt.Println("great test!")
			}
		})
	}

}

func TestRequestHandlerBadJson(t *testing.T) {
	testCasesFail := []struct {
		name           string
		expression     string
		expectedResult application.ServerErrorAnswer
	}{
		{
			name:           "int_value",
			expression:     `{"expression":9}`,
			expectedResult: application.ServerErrorAnswer{Expression: "", Error: "error in parsing json"},
		},
		{
			name:           "bool_value",
			expression:     `{"expression":true}`,
			expectedResult: application.ServerErrorAnswer{Expression: "", Error: "error in parsing json"},
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.expression)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", reader)
			w := httptest.NewRecorder()
			application.CalculationHandler(w, req)
			res := w.Result()
			defer res.Body.Close()

			var ans application.ServerErrorAnswer
			json.NewDecoder(res.Body).Decode(&ans)

			if ans != testCase.expectedResult {
				t.Fatalf("%+v\n should be equal %+v\n", ans, testCase.expectedResult)
			} else if res.StatusCode != 422 {
				t.Fatalf("%s test returns wrong status code", testCase.expression)
			} else {
				fmt.Println("great test!")
			}
		})
	}

}
