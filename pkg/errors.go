package calculator

import "errors"

var (
	ErrInvalidExpression  = errors.New("invalid expression")
	ErrDivisionByZero     = errors.New("division by zero")
	ErrRepeatingOperators = errors.New("repeating operators")
	ErrInvalidBrackets    = errors.New("invalid brackets")
	ErrInvalidOperator    = errors.New("invalid operator")
	ErrEmptyString        = errors.New("empty string")
)
