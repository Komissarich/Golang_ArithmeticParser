package calculator

import "errors"

var (
	ErrInvalidLetter      = errors.New("expression is not valid")
	ErrDivisionByZero     = errors.New("internal server error: division by zero")
	ErrRepeatingOperators = errors.New("internal server error: repeating operators")
	ErrInvalidBrackets    = errors.New("internal server error: invalid brackets")
	ErrInvalidOperator    = errors.New("internal server error: invalid operator")
	ErrEmptyString        = errors.New("internal server error: empty string")
)
