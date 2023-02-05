package expcontroller_test

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	expstore "github.com/gmaschi/log-exp-eval/internal/services/datastore/postgresql/exp"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

// eqCreateExpMatcher is a matcher type to validate the create expression method
type eqCreateExpMatcher struct {
	arg expstore.CreateExpressionParams
	id  uuid.UUID
}

func (e eqCreateExpMatcher) Matches(x interface{}) bool {
	arg, ok := x.(expstore.CreateExpressionParams)
	if !ok {
		return false
	}

	_, err := uuid.Parse(e.id.String())
	if err != nil {
		return false
	}
	e.arg.ExpressionID = arg.ExpressionID

	if math.Abs(float64(e.arg.CreatedAt.Second())-float64(arg.CreatedAt.Second())) > 1 ||
		math.Abs(float64(e.arg.UpdatedAt.Second())-float64(arg.UpdatedAt.Second())) > 1 {
		return false
	}
	e.arg.CreatedAt = arg.CreatedAt
	e.arg.UpdatedAt = arg.UpdatedAt

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateExpMatcher) String() string {
	return fmt.Sprintf("matches arg %v and uuid %v", e.arg, e.id)
}

func EqCreateExpParams(arg expstore.CreateExpressionParams, id uuid.UUID) gomock.Matcher {
	return eqCreateExpMatcher{arg, id}
}

// eqUpdateExpMatcher is a matcher type to validate the create expression method
type eqUpdateExpMatcher struct {
	arg expstore.UpdateExpressionParams
}

func (e eqUpdateExpMatcher) Matches(x interface{}) bool {
	arg, ok := x.(expstore.UpdateExpressionParams)
	if !ok {
		return false
	}

	if math.Abs(float64(e.arg.UpdatedAt.Second())-float64(arg.UpdatedAt.Second())) > 1 {
		return false
	}
	e.arg.UpdatedAt = arg.UpdatedAt

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateExpMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqUpdateExpParams(arg expstore.UpdateExpressionParams) gomock.Matcher {
	return eqUpdateExpMatcher{arg}
}
