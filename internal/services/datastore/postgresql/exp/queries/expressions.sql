-- name: GetExpressionByID :one
SELECT *
FROM expressions
WHERE expression_id = $1
    LIMIT 1;

-- name: ListExpressions :many
SELECT *
FROM expressions
ORDER BY row_id;

-- name: ListPaginatedExpressions :many
SELECT *
FROM expressions
ORDER BY row_id
    LIMIT $1 OFFSET $2;

-- name: CreateExpression :one
INSERT INTO expressions (expression_id, expression, username, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: UpdateExpression :one
UPDATE expressions
SET (expression, updated_at) = ($2, $3)
WHERE expression_id = $1
    RETURNING *;

-- name: DeleteExpressionByID :exec
DELETE
FROM expressions
WHERE expression_id = $1;