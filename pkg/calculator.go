package calculator

import (
	"slices"
	"strconv"
	"strings"
	"unicode"
)

func priority(char string) int {

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

func createPostfix(expression string) ([]string, error) {
	ind := 0
	num_ind := 0
	op_stack := []string{}
	res := []string{}
	if strings.Count(expression, "(") != strings.Count(expression, ")") {
		return nil, ErrInvalidBrackets
	}
	for {
		if ind >= len(expression) {
			break
		}

		if unicode.IsDigit(rune(expression[ind])) {
			number := string(expression[ind])
			num_ind = ind + 1
			for num_ind < len(expression) {
				if unicode.IsDigit(rune(expression[num_ind])) {
					number += string(expression[num_ind])
					num_ind += 1
				} else if expression[num_ind] == '.' {
					number += "."
					num_ind += 1
				} else {
					break
				}
			}

			res = append(res, string(number))
			ind = num_ind - 1
		}

		if isOperator(string(expression[ind])) {
			if ind != len(expression)-1 && isOperator(string(expression[ind])) && isOperator(string(expression[ind+1])) {
				return []string{}, ErrRepeatingOperators
			}
			if ind == len(expression)-1 {
				return []string{}, ErrInvalidOperator
			}
			for len(op_stack) > 0 && op_stack[len(op_stack)-1] != "(" && priority(string(expression[ind])) <= priority(op_stack[len(op_stack)-1]) {
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

		if expression[ind] != '(' && expression[ind] != ')' && !unicode.IsDigit(rune(expression[ind])) && !isOperator(string(expression[ind])) {
			return []string{}, ErrInvalidLetter
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
	postfix, err := createPostfix(expression)
	if err != nil {
		return 0, err
	}
	stack := []float64{}

	if len(postfix) == 0 {
		return 0, ErrEmptyString
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
					return 0, ErrDivisionByZero
				}
				stack = append(stack, sec_pop_item/fir_pop_item)
			}
		}

	}
	return stack[len(stack)-1], nil
}
