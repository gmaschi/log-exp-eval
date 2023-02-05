package eval_test

import (
	"github.com/gmaschi/log-exp-eval/internal/services/eval"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEvaluator_EvalLogicExp(t *testing.T) {
	evaluator := eval.New()

	testCases := []struct {
		name       string
		expression string
		expRes     bool
	}{
		{
			name:       "Simple true expression",
			expression: "1",
			expRes:     true,
		},
		{
			name:       "Simple true expression with operator",
			expression: "1 OR 0",
			expRes:     true,
		},
		{
			name:       "Nested true expression",
			expression: "((1 OR 0) AND (1 OR 0) OR 1)",
			expRes:     true,
		},
		// TODO: uncomment after operator precedence is implemented
		//{
		//	name:       "Nested true expression to evaluate operator precedence",
		//expression: "((0 OR 0) AND (1 OR 0) OR 1)",
		//expRes:     true,
		//},
		{
			name:       "Simple false expression",
			expression: "0",
			expRes:     false,
		},
		{
			name:       "Simple false expression with operator",
			expression: "(1 AND 0)",
			expRes:     false,
		},
		{
			name:       "Nested false expression",
			expression: "(((0 OR 0) AND (1 OR 0)) OR 0)",
			expRes:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := evaluator.EvalLogicExp(tc.expression)
			require.Equal(t, tc.expRes, res, "expected %q to be evaluated to %v", tc.expression, tc.expRes)
		})
	}
}
