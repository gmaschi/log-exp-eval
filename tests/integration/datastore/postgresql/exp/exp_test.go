package expstore_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	expstore "github.com/gmaschi/log-exp-eval/internal/services/datastore/postgresql/exp"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateExpression(t *testing.T) {
	createRandomExpression(t)
}

func TestGetExpression(t *testing.T) {
	t.Run("Get expression by ID", func(t *testing.T) {
		exp := createRandomExpression(t)

		gotExp, err := testStore.GetExpressionByID(context.Background(), exp.ExpressionID)
		require.NoError(t, err)
		require.NotEmpty(t, gotExp)
		require.Equal(t, exp.ExpressionID, gotExp.ExpressionID)
		require.Equal(t, exp.Expression, gotExp.Expression)
		require.Equal(t, exp.Username, gotExp.Username)
		require.WithinDuration(t, exp.CreatedAt, gotExp.CreatedAt, time.Second)
		require.WithinDuration(t, exp.UpdatedAt, gotExp.UpdatedAt, time.Second)
	})
}

func TestListExpressions(t *testing.T) {
	t.Run("List all expressions", func(t *testing.T) {
		n := 5
		for i := 0; i < n; i++ {
			createRandomExpression(t)
		}

		exps, err := testStore.ListExpressions(context.Background())
		require.NoError(t, err)
		require.NotEmpty(t, exps)
		require.GreaterOrEqual(t, len(exps), n)

		for _, exp := range exps {
			require.NotEmpty(t, exp)
		}
	})

	t.Run("List paginated expressions", func(t *testing.T) {
		n := 10
		for i := 0; i < n; i++ {
			createRandomExpression(t)
		}

		limit := 5
		offset := 3

		listPagArgs := expstore.ListPaginatedExpressionsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		}

		exps, err := testStore.ListPaginatedExpressions(context.Background(), listPagArgs)
		require.NoError(t, err)
		require.NotEmpty(t, exps)
		require.Len(t, exps, limit)

		for _, exp := range exps {
			require.NotEmpty(t, exp)
		}
	})
}

func TestUpdateExpression(t *testing.T) {
	t.Run("Update expression", func(t *testing.T) {
		exp := createRandomExpression(t)

		updateArgs := expstore.UpdateExpressionParams{
			ExpressionID: exp.ExpressionID,
			Expression:   "(h OR j)",
			UpdatedAt:    time.Now(),
		}

		updatedExp, err := testStore.UpdateExpression(context.Background(), updateArgs)
		require.NoError(t, err)
		require.NotEmpty(t, updatedExp)
		require.Equal(t, updateArgs.ExpressionID, updatedExp.ExpressionID)
		require.Equal(t, updateArgs.Expression, updatedExp.Expression)
		require.Equal(t, exp.Username, updatedExp.Username)
		require.WithinDuration(t, exp.CreatedAt, updatedExp.CreatedAt, time.Second)
		require.WithinDuration(t, updateArgs.UpdatedAt, updatedExp.UpdatedAt, time.Second)
	})
}

func TestDeleteExpression(t *testing.T) {
	t.Run("Delete expression by ID", func(t *testing.T) {
		exp := createRandomExpression(t)

		err := testStore.DeleteExpressionByID(context.Background(), exp.ExpressionID)
		require.NoError(t, err)

		deletedExp, err := testStore.GetExpressionByID(context.Background(), exp.ExpressionID)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, deletedExp)
	})
}

func createRandomExpression(t *testing.T) expstore.Expressions {
	t.Helper()

	expID, err := uuid.NewRandom()
	require.NoError(t, err)
	require.NotEmpty(t, expID)

	now := time.Now()
	createExpArgs := expstore.CreateExpressionParams{
		ExpressionID: expID,
		Expression:   "(x AND y) AND k",
		Username:     "some-username",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	exp, err := testStore.CreateExpression(context.Background(), createExpArgs)
	require.NoError(t, err)
	require.NotEmpty(t, exp)
	require.Equal(t, createExpArgs.ExpressionID, exp.ExpressionID)
	require.Equal(t, createExpArgs.Expression, exp.Expression)
	require.Equal(t, createExpArgs.Username, exp.Username)
	require.WithinDuration(t, createExpArgs.CreatedAt, exp.CreatedAt, time.Second)
	require.WithinDuration(t, createExpArgs.UpdatedAt, exp.UpdatedAt, time.Second)

	return exp
}
