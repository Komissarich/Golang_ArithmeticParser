package calculator

import (
	"calc/pkg/calculator/calc_errors"
	"slices"
	"strconv"
	"strings"
	"time"
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

func IsOperator(op string) bool {
	opers := []string{"-", "+", "*", "/", "^"}
	if slices.Contains(opers, op) {
		return true
	}
	return false
}

func LongAdd(a float64, b float64) float64 {
	time.Sleep(3 * time.Second)
	return a + b
}

func CreatePostfix(expression string) ([]string, error) {
	ind := 0
	num_ind := 0
	op_stack := []string{}
	res := []string{}
	if strings.Count(expression, "(") != strings.Count(expression, ")") {
		return nil, calc_errors.ErrInvalidBrackets
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

		if IsOperator(string(expression[ind])) {
			if ind != len(expression)-1 && IsOperator(string(expression[ind])) && IsOperator(string(expression[ind+1])) {
				return []string{}, calc_errors.ErrRepeatingOperators
			}
			if ind == len(expression)-1 {
				return []string{}, calc_errors.ErrInvalidOperator
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

		if expression[ind] != '(' && expression[ind] != ')' && !unicode.IsDigit(rune(expression[ind])) && !IsOperator(string(expression[ind])) {
			return []string{}, calc_errors.ErrInvalidLetter
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
		return 0, err
	}
	stack := []float64{}
	for _, val := range postfix {
		conv_val, err := strconv.ParseFloat(val, 64)
		if err == nil {
			stack = append(stack, conv_val)
		} else if IsOperator(val) {

			fir_pop_item := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			sec_pop_item := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			switch val {
			case "+":
				stack = append(stack, float64(LongAdd(fir_pop_item, sec_pop_item)))
			case "-":
				stack = append(stack, sec_pop_item-fir_pop_item)
			case "*":
				stack = append(stack, fir_pop_item*sec_pop_item)
			case "/":
				if fir_pop_item == 0 {
					return 0, calc_errors.ErrDivisionByZero
				}
				stack = append(stack, sec_pop_item/fir_pop_item)
			}
		}

	}
	return stack[len(stack)-1], nil
}
