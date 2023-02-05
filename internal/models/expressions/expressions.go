package expmodel

import (
	"github.com/google/uuid"
	"time"
)

type (
	// CreateExpressionRequest describes the request to create an expression.
	CreateExpressionRequest struct {
		Expression string `json:"expression" binding:"required"`
	}

	// CreateExpressionResponse describes the response when creating an expression.
	CreateExpressionResponse struct {
		RowID        int64     `json:"-"`
		ExpressionID uuid.UUID `json:"expressionID"`
		Expression   string    `json:"expression"`
		Username     string    `json:"username"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	// GetExpressionRequest describes the request to retrieve an expression.
	GetExpressionRequest struct {
		ID string `uri:"id" binding:"required"`
	}

	// GetExpressionResponse describes the response when getting an expression.
	GetExpressionResponse struct {
		RowID        int64     `json:"-"`
		ExpressionID uuid.UUID `json:"expressionID"`
		Expression   string    `json:"expression"`
		Username     string    `json:"username"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	// ListExpressionsRequest describes the request to list expressions.
	ListExpressionsRequest struct {
		PageID   int32 `form:"page_id"`
		PageSize int32 `form:"page_size"`
	}

	// ListExpressionsResponse describes the response when listing expressions.
	ListExpressionsResponse struct {
		RowID        int64     `json:"-"`
		ExpressionID uuid.UUID `json:"expressionID"`
		Expression   string    `json:"expression"`
		Username     string    `json:"username"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	// UpdateExpressionRequest describes the request to update an expression.
	UpdateExpressionRequest struct {
		ExpressionID string `json:"expression_id" binding:"required"`
		Expression   string `json:"expression"`
	}

	// UpdateExpressionResponse describes the response when updating an expression.
	UpdateExpressionResponse struct {
		RowID        int64     `json:"-"`
		ExpressionID uuid.UUID `json:"expressionID"`
		Expression   string    `json:"expression"`
		Username     string    `json:"username"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	// DeleteExpressionRequest describes the request to delete an expression.
	DeleteExpressionRequest struct {
		ID string `uri:"id" binding:"required"`
	}

	// EvaluateExpressionRequest describes the request to evaluate an expression.
	EvaluateExpressionRequest struct {
		ID string `uri:"id" binding:"required"`
	}

	// EvaluateExpressionResponse describes the request to evaluate an expression.
	EvaluateExpressionResponse struct {
		Result bool `json:"result"`
	}
)
