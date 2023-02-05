package eval

import "strings"

type (
	// Evaluator defines the method set to evaluate expressions.
	Evaluator interface {
		IsValidLogicExp(exp string) bool
		EvalLogicExp(exp string) bool
	}

	eval struct{}
)

// New instantiates a new evaluator
func New() Evaluator {
	return &eval{}
}

// IsValidLogicExp evaluates if a given logical expression is valid.
func (e *eval) IsValidLogicExp(exp string) bool {
	// TODO: implement me.
	return true
}

// EvalLogicExp evaluates a logical expression and returns the result of the expression.
func (e *eval) EvalLogicExp(exp string) bool {
	return evalExpression(exp)
}

// evalExpression is a helper function to evaluate the result of a logical expression.
// TODO: implement operator's order of precedence when evaluating an expression to cover some edge cases.
func evalExpression(exp string) bool {
	evalExp := strings.ReplaceAll(exp, " ", "")
	stack := make([]bool, 0)
	for i := 0; i < len(evalExp); i++ {
		switch evalExp[i] {
		case '1':
			stack = append(stack, true)
		case '0':
			stack = append(stack, false)
		case '(':
		case ')':
			val2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			val1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			operator := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			switch operator {
			case true:
				stack = append(stack, val1 || val2)
			case false:
				stack = append(stack, val1 && val2)
			}
		case 'A':
			if i+3 < len(evalExp) && evalExp[i:i+3] == "AND" {
				i += 2
				stack = append(stack, false)
			}
		case 'O':
			if i+2 < len(evalExp) && evalExp[i:i+2] == "OR" {
				i++
				stack = append(stack, true)
			}
		}
	}

	return stack[0]
}
