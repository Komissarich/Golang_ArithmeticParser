package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"unicode"
)

func Priority(char string) int {

	switch char {
	case "^":
		return 6
	case "*":
		return 5
	case "/":
		return 5
	case "%":
		return 3
	case "+":
		return 2
	case "-":
		return 2
	default:
		return 0
	}

}

func isOperator(op string) bool {
	opers := []string{"-", "+", "*", "/", "^"}
	if slices.Contains(opers, op) {
		return true
	}
	return false

}

func CreatePostfix(expression string) ([]string, error) {
	ind := 0
	op_stack := []string{}
	res := []string{}

	for {
		if ind >= len(expression) {
			break
		}

		if unicode.IsDigit(rune(expression[ind])) {
			number := string(expression[ind])
			if ind != len(expression)-1 {
				ind += 1

				for ind < len(expression) {
					if unicode.IsDigit(rune(expression[ind])) {
						number += string(expression[ind])
						ind += 1
					} else if expression[ind] == '.' {
						number += "."
						ind += 1
					} else {
						break
					}
				}
			}

			res = append(res, string(number))
		}
		if isOperator(string(expression[ind])) {
			if ind != len(expression)-1 && isOperator(string(expression[ind])) && isOperator(string(expression[ind+1])) {
				return []string{}, errors.New("bad string")
			}
			if ind == len(expression)-1 {
				return []string{}, errors.New("bad string")
			}
			for len(op_stack) > 0 && op_stack[len(op_stack)-1] != "(" && Priority(string(expression[ind])) <= Priority(op_stack[len(op_stack)-1]) {
				res = append(res, op_stack[len(op_stack)-1])
				op_stack = op_stack[:len(op_stack)-1]
			}
			op_stack = append(op_stack, string(expression[ind]))

		}
		if expression[ind] == ')' {

			for len(op_stack) > 0 {
				x := op_stack[len(op_stack)-1]
				op_stack = op_stack[:len(op_stack)-1]
				if x == "(" {
					break
				}
				res = append(res, x)
			}
		}
		if expression[ind] == '(' {
			op_stack = append(op_stack, string(expression[ind]))
		}
		ind += 1

	}
	for len(op_stack) > 0 {
		res = append(res, op_stack[len(op_stack)-1])
		op_stack = op_stack[:len(op_stack)-1]
	}
	return res, nil
}

func Calc(expression string) (float64, error) {
	postfix, err := CreatePostfix(expression)

	if err != nil {
		return 0, errors.New("something happened")
	}
	stack := []float64{}

	if len(postfix) == 0 {
		return 0, errors.New("Empty string")
	}
	for _, val := range postfix {
		conv_val, err := strconv.ParseFloat(val, 64)
		if err == nil {
			stack = append(stack, conv_val)
		} else if isOperator(val) {

			fir_pop_item := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			sec_pop_item := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			switch val {
			case "+":
				stack = append(stack, fir_pop_item+sec_pop_item)
			case "-":
				stack = append(stack, sec_pop_item-fir_pop_item)
			case "*":
				stack = append(stack, fir_pop_item*sec_pop_item)
			case "/":
				if fir_pop_item == 0 {
					return 0, errors.New("Division by zero")
				}
				stack = append(stack, sec_pop_item/fir_pop_item)
			}
		}

	}
	return stack[len(stack)-1], nil
}

func main() {
	str := "15/(7-(1+1))*3-(2+(1+1))*15/(7-(200+1))*3-(2+(1+1))*(15/(7-(1+1))*3-(2+(1+1))+15/(7-(1+1))*3-(2+(1+1)))"
	//right answer -30.072164948453608 <nil>
	fmt.Println(Calc(str))
}
