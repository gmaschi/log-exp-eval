package exp_docs

import (
	expmodel "github.com/gmaschi/log-exp-eval/internal/models/expressions"
	experrors "github.com/gmaschi/log-exp-eval/internal/models/expressions/errors"
)

// swagger:route POST /v1/expressions Expressions createExpressionParams
// Creates an expression and returns the created expression to the user.
//
// This route can only be used by authenticated users.
// responses:
//   201: createExpressionResponseWrapper
//   400: createExpressionBadRequest
//   401: createExpressionUnauthorized
//   500: createExpressionInternalServerError
//
//     Security:
//       bearer-normal:

// swagger:parameters createExpressionParams
type createExpressionParamsWrapper struct {
	// The request body contains the required information to create an expression.
	// in:body
	// required: true
	Body expmodel.CreateExpressionRequest
}

// The response body contains the information of the created expression.
// swagger:response
type createExpressionResponseWrapper struct {
	// in:body
	Body expmodel.CreateExpressionResponse
}

// Error response when the request body is not well formatted.
// swagger:response
type createExpressionBadRequest struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when the user does not provide authorization information to perform the request.
// swagger:response
type createExpressionUnauthorized struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when there is an internal server error.
// swagger:response
type createExpressionInternalServerError struct {
	// in:body
	Body experrors.ErrorResponse
}

// swagger:route GET /v1/expressions/{id} Expressions getExpressionParams
// Retrieves an expression.
//
// This route can only be used by authenticated users and a user can only retrieve expressions that he/she created.
// responses:
//   200: getExpressionResponseWrapper
//   400: getExpressionBadRequest
//   401: getExpressionUnauthorized
//   403: getExpressionForbidden
//   500: getExpressionInternalServerError
//
//     Security:
//       bearer-normal:

// swagger:parameters getExpressionParams
type getExpressionParamsWrapper struct {
	// The ID of the expression.
	// in:path
	ID string `json:"id"`
}

// The response body contains the information of the created expression.
// swagger:response
type getExpressionResponseWrapper struct {
	// in:body
	Body expmodel.CreateExpressionResponse
}

// Error response when the request body is not well formatted.
// swagger:response
type getExpressionBadRequest struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when the user does not provide authorization information to perform the request.
// swagger:response
type getExpressionUnauthorized struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when the user didn't create the expression he/she is trying to retrieve.
// swagger:response
type getExpressionForbidden struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when there is an internal server error.
// swagger:response
type getExpressionInternalServerError struct {
	// in:body
	Body experrors.ErrorResponse
}

// swagger:route GET /v1/expressions Expressions listExpressionsParams
// Retrieves a list of expressions.
//
// Pagination parameters shall be sent together to query paginated data or not sent at all to query all rows.
//
// This route can only be used by authenticated users.
// responses:
//   200: listExpressionsResponseWrapper
//   400: listExpressionsBadRequest
//   401: listExpressionsUnauthorized
//   500: listExpressionsInternalServerError
//
//     Security:
//       bearer-normal:

// swagger:parameters listExpressionsParams
type listExpressionsParamsWrapper struct {
	// The ID of the page to be returned. Min: 1.
	// in:query
	// min: 1
	PageID int64 `json:"page_id"`

	// The number of items per page to be retrieved. Min: 1.
	// in:query
	// min: 1
	PageSize int64 `json:"page_size"`
}

// The response body contains the information of the created expression.
// swagger:response
type listExpressionsResponseWrapper struct {
	// in:body
	Body []expmodel.CreateExpressionResponse
}

// Error response when the request body is not well formatted.
// swagger:response
type listExpressionsBadRequest struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when the user does not provide authorization information to perform the request.
// swagger:response
type listExpressionsUnauthorized struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when there is an internal server error.
// swagger:response
type listExpressionsInternalServerError struct {
	// in:body
	Body experrors.ErrorResponse
}

// swagger:route PATCH /v1/expressions Expressions updateExpressionParams
// Updates an expression and returns the updated expression to the user.
//
// This route can only be used by authenticated users and a user can only retrieve expressions that he/she created.
// responses:
//   200: updateExpressionResponseWrapper
//   400: updateExpressionBadRequest
//   401: updateExpressionUnauthorized
//   403: updateExpressionForbidden
//   500: updateExpressionInternalServerError
//
//     Security:
//       bearer-normal:

// swagger:parameters updateExpressionParams
type updateExpressionParamsWrapper struct {
	// The request body contains the information to update an expression.
	// in:body
	// required: true
	Body expmodel.UpdateExpressionRequest
}

// The response body contains the information of the created expression.
// swagger:response
type updateExpressionResponseWrapper struct {
	// in:body
	Body expmodel.UpdateExpressionResponse
}

// Error response when the request body is not well formatted.
// swagger:response
type updateExpressionBadRequest struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when the user does not provide authorization information to perform the request.
// swagger:response
type updateExpressionUnauthorized struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when the user does not have ownership of the expression.
// swagger:response
type updateExpressionForbidden struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when there is an internal server error.
// swagger:response
type updateExpressionInternalServerError struct {
	// in:body
	Body experrors.ErrorResponse
}

// swagger:route DELETE /v1/expressions/{id} Expressions deleteExpressionParams
// Deletes an expression.
//
// This route can only be used by authenticated users and a user can only delete expressions that he/she created.
// responses:
//   204: deleteExpressionResponseWrapper
//   400: deleteExpressionBadRequest
//   401: deleteExpressionUnauthorized
//   403: deleteExpressionForbidden
//   500: deleteExpressionInternalServerError
//
//     Security:
//       bearer-normal:

// swagger:parameters deleteExpressionParams
type deleteExpressionParamsWrapper struct {
	// The ID of the expression.
	// in:path
	ID string `json:"id"`
}

// The response body contains the information of the created expression.
// swagger:response
type deleteExpressionResponseWrapper struct {
	// in:body
}

// Error response when the request body is not well formatted.
// swagger:response
type deleteExpressionBadRequest struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when the user does not provide authorization information to perform the request.
// swagger:response
type deleteExpressionUnauthorized struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when the user didn't create the expression he/she is trying to retrieve.
// swagger:response
type deleteExpressionForbidden struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when there is an internal server error.
// swagger:response
type deleteExpressionInternalServerError struct {
	// in:body
	Body experrors.ErrorResponse
}

// swagger:route GET /v1/evaluate/{id} Expressions evaluateExpressionParams
// Evaluates an expression and returns its value.
//
// All parameters required by the given expression must be sent as query parameters to perform the evaluation. If at least one parameter is missing, an error will be returned to the user.
//
// This route can only be used by authenticated users.
// responses:
//   200: evaluateExpressionResponseWrapper
//   400: evaluateExpressionBadRequest
//   401: evaluateExpressionUnauthorized
//   500: evaluateExpressionInternalServerError
//
//     Security:
//       bearer-normal:

// swagger:parameters evaluateExpressionParams
type evaluateExpressionParamsWrapper struct {
	// The ID of the expression.
	// in:path
	ID string `json:"id"`

	// Some X variable parameter
	// in:query
	X string `json:"x"`

	// Some Y variable parameter
	// in:query
	Y string `json:"y"`
}

// The response body contains the information of the created expression.
// swagger:response
type evaluateExpressionResponseWrapper struct {
	// in:body
	Body expmodel.CreateExpressionResponse
}

// Error response when the request body is not well formatted.
// swagger:response
type evaluateExpressionBadRequest struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when the user does not provide authorization information to perform the request.
// swagger:response
type evaluateExpressionUnauthorized struct {
	// in:body
	Body experrors.ErrorResponse
}

// Error response when there is an internal server error.
// swagger:response
type evaluateExpressionInternalServerError struct {
	// in:body
	Body experrors.ErrorResponse
}
