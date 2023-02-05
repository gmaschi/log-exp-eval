package experrors

import "errors"

type (
	ExpressionError string

	ErrorResponse struct {
		Error string `json:"error"`
	}
)

const (
	// general
	ErrInvalidExpressionID       ExpressionError = "invalid expression id"
	ErrInvalidExpression         ExpressionError = "invalid expression"
	ErrInvalidArguments          ExpressionError = "invalid arguments"
	ErrCreatingExpression        ExpressionError = "failed to create expression"
	ErrRetrievingExpression      ExpressionError = "failed to retrieve expression"
	ErrUpdatingExpression        ExpressionError = "failed to update expression"
	ErrDeletingExpression        ExpressionError = "failed to delete expression"
	ErrInvalidEvaluateExpression ExpressionError = "failed to evaluate expression"

	ErrRecordNotFound ExpressionError = "record not found"
	ErrInternalServer ExpressionError = "internal error"
	ErrBadRequest     ExpressionError = "bad request"

	ErrRequiredPaginationData ExpressionError = "list paginated requires both page_id and page_size"
)

func (se ExpressionError) Error() error {
	return errors.New(se.String())
}

func (se ExpressionError) String() string {
	return string(se)
}
