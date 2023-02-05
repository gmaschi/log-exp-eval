package expcontroller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	authmid "github.com/gmaschi/log-exp-eval/internal/controllers/middlewares/auth-mid"
	expmodel "github.com/gmaschi/log-exp-eval/internal/models/expressions"
	experrors "github.com/gmaschi/log-exp-eval/internal/models/expressions/errors"
	expstore "github.com/gmaschi/log-exp-eval/internal/services/datastore/postgresql/exp"
	"github.com/gmaschi/log-exp-eval/internal/services/eval"
	"github.com/gmaschi/log-exp-eval/pkg/tools/marshaller"
	ginmidctx "github.com/gmaschi/log-exp-eval/pkg/tools/middlewares/gin/context"
	"github.com/gmaschi/log-exp-eval/pkg/tools/pagination"
	"github.com/gmaschi/log-exp-eval/pkg/tools/parse"
	"github.com/gmaschi/log-exp-eval/pkg/tools/str"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Controller defines the expression controllers and its required fields.
type Controller struct {
	store     expstore.Store
	evaluator eval.Evaluator
}

// New creates a pointer to a Controller
func New(store expstore.Store, ev eval.Evaluator) *Controller {
	return &Controller{
		store:     store,
		evaluator: ev,
	}
}

// Create handles the request to create an expression.
func (c *Controller) Create(ctx *gin.Context) {
	var req expmodel.CreateExpressionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, parse.ErrorAsJSON(err))
		return
	}

	// TODO: implement expression validator
	exp := strings.TrimSpace(req.Expression)
	if exp == "" {
		ctx.JSON(
			http.StatusBadRequest,
			parse.ErrorAsJSON(fmt.Errorf("%s: empty expression", experrors.ErrInvalidExpressionID)),
		)
		return
	}

	expID, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			parse.ErrorAsJSON(fmt.Errorf("failed to generate expression ID")),
		)
		return
	}

	authPayload, err := extractAuthPayload(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			parse.ErrorAsJSON(err),
		)
		return
	}

	now := time.Now()
	createExpArgs := expstore.CreateExpressionParams{
		ExpressionID: expID,
		Expression:   req.Expression,
		Username:     authPayload.Username,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	createdExp, err := c.store.CreateExpression(ctx, createExpArgs)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			parse.ErrorAsJSON(experrors.ErrCreatingExpression.Error()),
		)
		return
	}

	var res expmodel.CreateExpressionResponse
	err = marshaller.Response(createdExp, &res)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			parse.ErrorAsJSON(fmt.Errorf("%s: %v", experrors.ErrInternalServer, err)),
		)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

// Get handles the request to get an expression by ID.
func (c *Controller) Get(ctx *gin.Context) {
	var req expmodel.GetExpressionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, parse.ErrorAsJSON(err))
		return
	}

	authPayload, err := extractAuthPayload(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			parse.ErrorAsJSON(err),
		)
		return
	}

	sanitizedID := strings.TrimSpace(req.ID)
	expID, err := uuid.Parse(sanitizedID)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			parse.ErrorAsJSON(fmt.Errorf("%s: %s", experrors.ErrInvalidExpressionID, sanitizedID)),
		)
		return
	}

	gotExp, err := c.store.GetExpressionByID(ctx, expID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(
				http.StatusNotFound,
				parse.ErrorAsJSON(experrors.ErrRecordNotFound.Error()),
			)
			return
		}

		ctx.JSON(http.StatusInternalServerError, parse.ErrorAsJSON(experrors.ErrRetrievingExpression.Error()))
		return
	}

	if gotExp.Username != authPayload.Username {
		ctx.JSON(
			http.StatusNotFound,
			parse.ErrorAsJSON(fmt.Errorf("expression  %s not found for user %s (%s)", req.ID, authPayload.Username, authPayload.UserID)),
		)
		return
	}

	var res expmodel.GetExpressionResponse
	err = marshaller.Response(gotExp, &res)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			parse.ErrorAsJSON(fmt.Errorf("%s: %v", experrors.ErrInternalServer, err)),
		)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// Delete handles the request to soft deletes an expression by ID
func (c *Controller) Delete(ctx *gin.Context) {
	var req expmodel.DeleteExpressionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, parse.ErrorAsJSON(err))
		return
	}

	authPayload, err := extractAuthPayload(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			parse.ErrorAsJSON(err),
		)
		return
	}

	sanitizedID := strings.TrimSpace(req.ID)
	expID, err := uuid.Parse(sanitizedID)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			parse.ErrorAsJSON(fmt.Errorf("%s: %s", experrors.ErrInvalidExpressionID, sanitizedID)),
		)
		return
	}

	gotExp, err := c.store.GetExpressionByID(ctx, expID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(
				http.StatusNotFound,
				parse.ErrorAsJSON(experrors.ErrRecordNotFound.Error()),
			)
			return
		}

		ctx.JSON(http.StatusInternalServerError, parse.ErrorAsJSON(experrors.ErrDeletingExpression.Error()))
		return
	}

	if gotExp.Username != authPayload.Username {
		ctx.JSON(
			http.StatusForbidden,
			parse.ErrorAsJSON(fmt.Errorf("expression %s does not belong to user %s (%s)", gotExp.ExpressionID, authPayload.Username, authPayload.UserID)),
		)
		return
	}

	err = c.store.DeleteExpressionByID(ctx, gotExp.ExpressionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, parse.ErrorAsJSON(experrors.ErrDeletingExpression.Error()))
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

// Update handles the request to update an expression by ID
func (c *Controller) Update(ctx *gin.Context) {
	var req expmodel.UpdateExpressionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, parse.ErrorAsJSON(err))
		return
	}

	authPayload, err := extractAuthPayload(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			parse.ErrorAsJSON(err),
		)
		return
	}

	sanitizedExp := strings.TrimSpace(req.Expression)
	if sanitizedExp == "" {
		ctx.JSON(http.StatusBadRequest, parse.ErrorAsJSON(fmt.Errorf("%s: empty expression", experrors.ErrInvalidExpression)))
		return
	}

	sanitizedID := strings.TrimSpace(req.ExpressionID)
	expID, err := uuid.Parse(sanitizedID)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			parse.ErrorAsJSON(fmt.Errorf("%s: %s", experrors.ErrInvalidExpressionID, sanitizedID)),
		)
		return
	}

	gotExp, err := c.store.GetExpressionByID(ctx, expID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(
				http.StatusNotFound,
				parse.ErrorAsJSON(experrors.ErrRecordNotFound.Error()),
			)
			return
		}

		ctx.JSON(http.StatusInternalServerError, parse.ErrorAsJSON(experrors.ErrRetrievingExpression.Error()))
		return
	}

	if gotExp.Username != authPayload.Username {
		ctx.JSON(
			http.StatusForbidden,
			parse.ErrorAsJSON(fmt.Errorf("expression  %s not found for user %s (%s)", expID, authPayload.Username, authPayload.UserID)),
		)
		return
	}

	updateExpArgs := expstore.UpdateExpressionParams{
		ExpressionID: gotExp.ExpressionID,
		Expression:   sanitizedExp,
		UpdatedAt:    time.Now(),
	}
	updatedExp, err := c.store.UpdateExpression(ctx, updateExpArgs)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(
				http.StatusNotFound,
				parse.ErrorAsJSON(experrors.ErrRecordNotFound.Error()),
			)
			return
		}

		ctx.JSON(http.StatusInternalServerError, parse.ErrorAsJSON(experrors.ErrUpdatingExpression.Error()))
		return
	}

	var res expmodel.UpdateExpressionResponse
	err = marshaller.Response(updatedExp, &res)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			parse.ErrorAsJSON(fmt.Errorf("%s: %v", experrors.ErrInternalServer, err)),
		)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// List handles the request to list expressions.
func (c *Controller) List(ctx *gin.Context) {
	var req expmodel.ListExpressionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			parse.ErrorAsJSON(fmt.Errorf("%s: %v", experrors.ErrBadRequest, err)),
		)
		return
	}

	var isPaginated bool
	pageID := req.PageID
	pageSize := req.PageSize
	switch {
	case pageID == 0 && pageSize == 0:
	case pageID > 0 && pageSize > 0:
		isPaginated = true
	default:
		ctx.JSON(
			http.StatusBadRequest,
			parse.ErrorAsJSON(experrors.ErrRequiredPaginationData.Error()),
		)
		return
	}

	exps, err := c.listExpressions(ctx, isPaginated, pageID, pageSize)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, []expstore.Expressions{})
			return
		}

		ctx.JSON(http.StatusInternalServerError, parse.ErrorAsJSON(experrors.ErrRecordNotFound.Error()))
		return
	}

	var res []expmodel.ListExpressionsResponse
	err = marshaller.Response(exps, &res)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			parse.ErrorAsJSON(fmt.Errorf("%s: %v", experrors.ErrInternalServer, err)),
		)
		return
	}

	ctx.JSON(http.StatusOK, exps)
}

// Evaluate handles the request to evaluate an expression.
func (c *Controller) Evaluate(ctx *gin.Context) {
	var req expmodel.EvaluateExpressionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, parse.ErrorAsJSON(err))
		return
	}

	sanitizedID := strings.TrimSpace(req.ID)
	expID, err := uuid.Parse(sanitizedID)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			parse.ErrorAsJSON(fmt.Errorf("%s: %s", experrors.ErrInvalidExpressionID, sanitizedID)),
		)
		return
	}

	gotExp, err := c.store.GetExpressionByID(ctx, expID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(
				http.StatusNotFound,
				parse.ErrorAsJSON(experrors.ErrRecordNotFound.Error()),
			)
			return
		}

		ctx.JSON(http.StatusInternalServerError, parse.ErrorAsJSON(experrors.ErrDeletingExpression.Error()))
		return
	}

	evalExp, err := c.getEvalExp(ctx, gotExp.Expression)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, parse.ErrorAsJSON(fmt.Errorf("%s: %v", experrors.ErrInvalidEvaluateExpression.Error(), err)))
		return
	}

	expResult := c.evaluator.EvalLogicExp(evalExp)
	res := expmodel.EvaluateExpressionResponse{
		Result: expResult,
	}

	ctx.JSON(http.StatusOK, res)
}

// getEvalExp is a helper function to replace the given query parameter values for the given expression
func (c *Controller) getEvalExp(ctx *gin.Context, exp string) (string, error) {
	qV := ctx.Request.URL.Query()
	evalExp := exp
	for k, v := range qV {
		if k != "" {
			// ignore multiple values
			fV := v[0]
			if fV != "0" && fV != "1" {
				return "", fmt.Errorf("variable values must be either 0 or 1. received %s: %s", k, fV)
			}

			lowerK := strings.ToLower(k)
			evalExp = strings.ReplaceAll(evalExp, lowerK, fV)
		}
	}

	lowers, hasLowers := str.HasLowers(evalExp)
	if hasLowers {
		return "", fmt.Errorf("missing required arguments: %+v", lowers)
	}

	return evalExp, nil
}

func (c *Controller) listExpressions(
	ctx *gin.Context,
	isPaginated bool,
	pageID,
	pageSize int32,
) ([]expstore.Expressions, error) {
	totalExps, err := c.store.ListExpressions(ctx)
	if err != nil {
		return nil, err
	}

	if isPaginated {
		listPagArgs := expstore.ListPaginatedExpressionsParams{
			Limit:  pageSize,
			Offset: pageSize * (pageID - 1),
		}

		pgData := pagination.GetData(len(totalExps), pageID, pageSize)
		ginmidctx.SetPaginationHeader(ctx, pgData)

		return c.store.ListPaginatedExpressions(ctx, listPagArgs)
	}

	return totalExps, nil
}

// extractAuthPayload extracts the payload information from the provided context.
func extractAuthPayload(ctx *gin.Context) (authmid.AuthValue, error) {
	payloadRawContent, ok := ctx.Get(authmid.AuthorizationPayloadKey)
	if !ok {
		return authmid.AuthValue{}, fmt.Errorf("authorization payload key not found")
	}

	authPayload, ok := payloadRawContent.(authmid.AuthValue)
	if !ok {
		return authmid.AuthValue{}, fmt.Errorf("wrong format provided for authorization payload: %v", payloadRawContent)
	}

	return authPayload, nil
}
