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
		expectedResult application.ServerAnswer
	}{
		{
			name:           "simple",
			expression:     `{"expression":"1+1"}`,
			expectedResult: application.ServerAnswer{Expression: "1+1", Result: 2, Error: ""},
		},
		{
			name:           "priority",
			expression:     `{"expression":"(2+2)*2"}`,
			expectedResult: application.ServerAnswer{Expression: "(2+2)*2", Result: 8, Error: ""},
		},
		{
			name:           "priority2",
			expression:     `{"expression":"2+2-2"}`,
			expectedResult: application.ServerAnswer{Expression: "2+2-2", Result: 2, Error: ""},
		},
		{
			name:           "division",
			expression:     `{"expression":"1/2"}`,
			expectedResult: application.ServerAnswer{Expression: "1/2", Result: 0.5, Error: ""},
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.expression)
			req := httptest.NewRequest(http.MethodGet, "/", reader)
			w := httptest.NewRecorder()
			application.CalculationHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			var ans application.ServerAnswer
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
		expectedResult application.ServerAnswer
	}{
		{
			name:           "operators",
			expression:     `{"expression":"1+*1"}`,
			expectedResult: application.ServerAnswer{Expression: "1+*1", Result: 0, Error: "repeating operators"},
		},
		{
			name:           "operator2",
			expression:     `{"expression":"2+2**2"}`,
			expectedResult: application.ServerAnswer{Expression: "2+2**2", Result: 0, Error: "repeating operators"},
		},
		{
			name:           "brackets",
			expression:     `{"expression":"((2+2-*(2"}`,
			expectedResult: application.ServerAnswer{Expression: "((2+2-*(2", Result: 0, Error: "invalid brackets"},
		},
		{
			name:           "empty",
			expression:     `{"expression":""}`,
			expectedResult: application.ServerAnswer{Expression: "", Result: 0, Error: "empty string"},
		},
		{
			name:           "division by zero",
			expression:     `{"expression":"2/0"}`,
			expectedResult: application.ServerAnswer{Expression: "2/0", Result: 0, Error: "division by zero"},
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.expression)
			req := httptest.NewRequest(http.MethodGet, "/", reader)
			w := httptest.NewRecorder()
			application.CalculationHandler(w, req)
			res := w.Result()
			defer res.Body.Close()

			var ans application.ServerAnswer
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
		expectedResult application.ServerAnswer
	}{
		{
			name:           "int_value",
			expression:     `{"expression":9}`,
			expectedResult: application.ServerAnswer{Expression: "", Result: 0, Error: "error in parsing json"},
		},
		{
			name:           "bool_value",
			expression:     `{"expression":true}`,
			expectedResult: application.ServerAnswer{Expression: "", Result: 0, Error: "error in parsing json"},
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.expression)
			req := httptest.NewRequest(http.MethodGet, "/", reader)
			w := httptest.NewRecorder()
			application.CalculationHandler(w, req)
			res := w.Result()
			defer res.Body.Close()

			var ans application.ServerAnswer
			json.NewDecoder(res.Body).Decode(&ans)

			if ans != testCase.expectedResult {
				t.Fatalf("%+v\n should be equal %+v\n", ans, testCase.expectedResult)
			} else if res.StatusCode != 400 {
				t.Fatalf("%s test returns wrong status code", testCase.expression)
			} else {
				fmt.Println("great test!")
			}
		})
	}

}
