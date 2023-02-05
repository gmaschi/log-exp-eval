package expcontroller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	authmid "github.com/gmaschi/log-exp-eval/internal/controllers/middlewares/auth-mid"
	expmodel "github.com/gmaschi/log-exp-eval/internal/models/expressions"
	expserver "github.com/gmaschi/log-exp-eval/internal/servers/expressions"
	expstore "github.com/gmaschi/log-exp-eval/internal/services/datastore/postgresql/exp"
	mockedexpstore "github.com/gmaschi/log-exp-eval/internal/services/datastore/postgresql/exp/mocks"
	mockedeval "github.com/gmaschi/log-exp-eval/internal/services/eval/mocks"
	"github.com/gmaschi/log-exp-eval/pkg/tools/config/env"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	authValue := authmid.AuthValue{
		BearerToken: authmid.BearerToken1,
		UserID:      "12345",
		Username:    "John Doe",
	}

	exp := getExp(t, authValue.Username)

	testCases := []struct {
		name          string
		body          map[string]interface{}
		bearerToken   authmid.BearerToken
		setupAuth     func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken)
		buildStubs    func(store *mockedexpstore.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Happy path",
			body: map[string]interface{}{
				"expression": exp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				createArg := expstore.CreateExpressionParams{
					ExpressionID: exp.ExpressionID,
					Expression:   exp.Expression,
					Username:     exp.Username,
					CreatedAt:    exp.CreatedAt,
					UpdatedAt:    exp.UpdatedAt,
				}
				store.EXPECT().
					CreateExpression(gomock.Any(), EqCreateExpParams(createArg, exp.ExpressionID)).
					Times(1).Return(exp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchCreate(t, recorder.Body, exp)
			},
		},
		{
			name: "Error - empty expression provided",
			body: map[string]interface{}{
				"expression": "",
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().CreateExpression(gomock.Any(), gomock.Any()).Times(0).Return(exp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Error - blank expression provided",
			body: map[string]interface{}{
				"expression": "  ",
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().CreateExpression(gomock.Any(), gomock.Any()).Times(0).Return(exp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Error - database - failed to create expression",
			body: map[string]interface{}{
				"expression": exp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				createArg := expstore.CreateExpressionParams{
					ExpressionID: exp.ExpressionID,
					Expression:   exp.Expression,
					Username:     exp.Username,
					CreatedAt:    exp.CreatedAt,
					UpdatedAt:    exp.UpdatedAt,
				}
				store.EXPECT().
					CreateExpression(gomock.Any(), EqCreateExpParams(createArg, exp.ExpressionID)).
					Times(1).Return(expstore.Expressions{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Error - no authorization",
			body: map[string]interface{}{
				"expression": exp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth:   func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().CreateExpression(gomock.Any(), gomock.Any()).Times(0).Return(exp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockedexpstore.NewMockStore(ctrl)
			tc.buildStubs(store)

			config, err := env.NewConfig()
			require.NoError(t, err)

			server, err := expserver.New(config, store, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := "/v1/expressions"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, tc.bearerToken)
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGet(t *testing.T) {
	authValue := authmid.AuthValue{
		BearerToken: authmid.BearerToken1,
		UserID:      "12345",
		Username:    "John Doe",
	}
	anotherAuthValue := authmid.AuthValue{
		BearerToken: authmid.BearerToken2,
		UserID:      "98765",
		Username:    "Jane Doe",
	}

	exp := getExp(t, authValue.Username)

	testCases := []struct {
		name          string
		id            string
		bearerToken   authmid.BearerToken
		setupAuth     func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken)
		buildStubs    func(store *mockedexpstore.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "Happy path",
			id:          exp.ExpressionID.String(),
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGet(t, recorder.Body, exp)
			},
		},
		{
			name:        "Error - invalid expression ID",
			id:          "some-invalid-id",
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(expstore.Expressions{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Error - database - expression not found",
			id:          exp.ExpressionID.String(),
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().
					GetExpressionByID(gomock.Any(), exp.ExpressionID).
					Times(1).Return(expstore.Expressions{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "Error - database - internal error",
			id:          exp.ExpressionID.String(),
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().
					GetExpressionByID(gomock.Any(), exp.ExpressionID).
					Times(1).Return(expstore.Expressions{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "Error - expression does not belong to authenticated user",
			id:          exp.ExpressionID.String(),
			bearerToken: anotherAuthValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "Error - not authenticated",
			id:          exp.ExpressionID.String(),
			bearerToken: authValue.BearerToken,
			setupAuth:   func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(expstore.Expressions{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockedexpstore.NewMockStore(ctrl)
			tc.buildStubs(store)

			config, err := env.NewConfig()
			require.NoError(t, err)

			server, err := expserver.New(config, store, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := "/v1/expressions/" + tc.id

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, tc.bearerToken)
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDelete(t *testing.T) {
	authValue := authmid.AuthValue{
		BearerToken: authmid.BearerToken1,
		UserID:      "12345",
		Username:    "John Doe",
	}
	anotherAuthValue := authmid.AuthValue{
		BearerToken: authmid.BearerToken2,
		UserID:      "98765",
		Username:    "Jane Doe",
	}

	exp := getExp(t, authValue.Username)

	testCases := []struct {
		name          string
		id            string
		bearerToken   authmid.BearerToken
		setupAuth     func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken)
		buildStubs    func(store *mockedexpstore.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "Happy path",
			id:          exp.ExpressionID.String(),
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)
				store.EXPECT().DeleteExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:        "Error - invalid expression ID",
			id:          "some-invalid-id",
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(expstore.Expressions{}, nil)
				store.EXPECT().DeleteExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Error - database - get expression not found",
			id:          exp.ExpressionID.String(),
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(expstore.Expressions{}, sql.ErrNoRows)
				store.EXPECT().DeleteExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "Error - database - get expression internal error",
			id:          exp.ExpressionID.String(),
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(expstore.Expressions{}, sql.ErrConnDone)
				store.EXPECT().DeleteExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "Error - expression does not belong to authenticated user",
			id:          exp.ExpressionID.String(),
			bearerToken: anotherAuthValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)
				store.EXPECT().DeleteExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:        "Error - database - delete expression error",
			id:          exp.ExpressionID.String(),
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)
				store.EXPECT().DeleteExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "Error - not authenticated",
			id:          exp.ExpressionID.String(),
			bearerToken: authValue.BearerToken,
			setupAuth:   func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(expstore.Expressions{}, nil)
				store.EXPECT().DeleteExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockedexpstore.NewMockStore(ctrl)
			tc.buildStubs(store)

			config, err := env.NewConfig()
			require.NoError(t, err)

			server, err := expserver.New(config, store, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := "/v1/expressions/" + tc.id

			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, tc.bearerToken)
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestList(t *testing.T) {
	authValue := authmid.AuthValue{
		BearerToken: authmid.BearerToken1,
		UserID:      "12345",
		Username:    "John Doe",
	}

	n := 10
	exps := make([]expstore.Expressions, 0, n)
	var exp expstore.Expressions
	for i := 0; i < n; i++ {
		exp = getExp(t, authValue.Username)
		exps = append(exps, exp)
	}
	pageID := 1
	pageSize := 5

	testCases := []struct {
		name          string
		queries       map[string]string
		bearerToken   authmid.BearerToken
		setupAuth     func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken)
		buildStubs    func(store *mockedexpstore.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "Happy path - list all",
			queries:     map[string]string{},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().ListExpressions(gomock.Any()).Times(1).Return(exps, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchList(t, recorder.Body, exps)
			},
		},
		{
			name: "Happy path - list paginated",
			queries: map[string]string{
				"page_id":   strconv.Itoa(pageID),
				"page_size": strconv.Itoa(pageSize),
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().ListExpressions(gomock.Any()).Times(1).Return(exps, nil)

				listPagArgs := expstore.ListPaginatedExpressionsParams{
					Limit:  int32(pageSize),
					Offset: int32(pageSize) * int32(pageID-1),
				}
				store.EXPECT().ListPaginatedExpressions(gomock.Any(), listPagArgs).Times(1).Return(exps[:5], nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchList(t, recorder.Body, exps[:5])
			},
		},
		{
			name: "Error - invalid pagination data - invalid page ID",
			queries: map[string]string{
				"page_id":   "invalid-page-id",
				"page_size": strconv.Itoa(pageSize),
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().ListExpressions(gomock.Any()).Times(0).Return(nil, nil)
				store.EXPECT().ListPaginatedExpressions(gomock.Any(), gomock.Any()).Times(0).Return(nil, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Error - invalid pagination data - invalid page size",
			queries: map[string]string{
				"page_id":   strconv.Itoa(pageID),
				"page_size": "invalid-page-size",
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().ListExpressions(gomock.Any()).Times(0).Return(nil, nil)
				store.EXPECT().ListPaginatedExpressions(gomock.Any(), gomock.Any()).Times(0).Return(nil, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Error - missing pagination data - missing page id",
			queries: map[string]string{
				"page_id":   "",
				"page_size": strconv.Itoa(pageSize),
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().ListExpressions(gomock.Any()).Times(0).Return(nil, nil)
				store.EXPECT().ListPaginatedExpressions(gomock.Any(), gomock.Any()).Times(0).Return(nil, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Error - missing pagination data - missing page size",
			queries: map[string]string{
				"page_id":   strconv.Itoa(pageID),
				"page_size": "",
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().ListExpressions(gomock.Any()).Times(0).Return(nil, nil)
				store.EXPECT().ListPaginatedExpressions(gomock.Any(), gomock.Any()).Times(0).Return(nil, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Error - failed to list expressions - not found",
			queries:     map[string]string{},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().ListExpressions(gomock.Any()).Times(1).Return([]expstore.Expressions{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "Error - failed to list expressions - internal error",
			queries:     map[string]string{},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().ListExpressions(gomock.Any()).Times(1).Return([]expstore.Expressions{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		//{
		//	name:        "Error - invalid expression ID",
		//	id:          "some-invalid-id",
		//	bearerToken: authValue.BearerToken,
		//	setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
		//		addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
		//	},
		//	buildStubs: func(store *mockedexpstore.MockStore) {
		//		store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(expstore.Expressions{}, nil)
		//	},
		//	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusBadRequest, recorder.Code)
		//	},
		//},
		//{
		//	name:        "Error - database - expression not found",
		//	id:          exp.ExpressionID.String(),
		//	bearerToken: authValue.BearerToken,
		//	setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
		//		addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
		//	},
		//	buildStubs: func(store *mockedexpstore.MockStore) {
		//		store.EXPECT().
		//			GetExpressionByID(gomock.Any(), exp.ExpressionID).
		//			Times(1).Return(expstore.Expressions{}, sql.ErrNoRows)
		//	},
		//	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusNotFound, recorder.Code)
		//	},
		//},
		//{
		//	name:        "Error - database - internal error",
		//	id:          exp.ExpressionID.String(),
		//	bearerToken: authValue.BearerToken,
		//	setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
		//		addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
		//	},
		//	buildStubs: func(store *mockedexpstore.MockStore) {
		//		store.EXPECT().
		//			GetExpressionByID(gomock.Any(), exp.ExpressionID).
		//			Times(1).Return(expstore.Expressions{}, sql.ErrConnDone)
		//	},
		//	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		//	},
		//},
		//{
		//	name:        "Error - expression does not belong to authenticated user",
		//	id:          exp.ExpressionID.String(),
		//	bearerToken: anotherAuthValue.BearerToken,
		//	setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
		//		addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
		//	},
		//	buildStubs: func(store *mockedexpstore.MockStore) {
		//		store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)
		//	},
		//	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusNotFound, recorder.Code)
		//	},
		//},
		//{
		//	name:        "Error - not authenticated",
		//	id:          exp.ExpressionID.String(),
		//	bearerToken: authValue.BearerToken,
		//	setupAuth:   func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {},
		//	buildStubs: func(store *mockedexpstore.MockStore) {
		//		store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(expstore.Expressions{}, nil)
		//	},
		//	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		//	},
		//},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockedexpstore.NewMockStore(ctrl)
			tc.buildStubs(store)

			config, err := env.NewConfig()
			require.NoError(t, err)

			server, err := expserver.New(config, store, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := "/v1/expressions"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := req.URL.Query()
			for k, v := range tc.queries {
				if v != "" {
					q.Set(k, v)
				}
			}
			req.URL.RawQuery = q.Encode()

			tc.setupAuth(t, req, tc.bearerToken)
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdate(t *testing.T) {
	authValue := authmid.AuthValue{
		BearerToken: authmid.BearerToken1,
		UserID:      "12345",
		Username:    "John Doe",
	}
	anotherAuthValue := authmid.AuthValue{
		BearerToken: authmid.BearerToken2,
		UserID:      "98765",
		Username:    "Jane Doe",
	}

	exp := getExp(t, authValue.Username)
	now := time.Now()
	updatedExp := expstore.Expressions{
		RowID:        exp.RowID,
		ExpressionID: exp.ExpressionID,
		Expression:   "(0 OR 0)",
		Username:     exp.Username,
		CreatedAt:    exp.CreatedAt,
		UpdatedAt:    now,
	}

	testCases := []struct {
		name          string
		body          map[string]interface{}
		bearerToken   authmid.BearerToken
		setupAuth     func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken)
		buildStubs    func(store *mockedexpstore.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Happy path",
			body: map[string]interface{}{
				"expression_id": updatedExp.ExpressionID,
				"expression":    updatedExp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)

				updateArg := expstore.UpdateExpressionParams{
					ExpressionID: updatedExp.ExpressionID,
					Expression:   updatedExp.Expression,
					UpdatedAt:    updatedExp.UpdatedAt,
				}
				store.EXPECT().
					UpdateExpression(gomock.Any(), EqUpdateExpParams(updateArg)).
					Times(1).Return(updatedExp, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUpdate(t, recorder.Body, updatedExp)
			},
		},
		{
			name: "Error - missing expression ID",
			body: map[string]interface{}{
				"expression_id": "",
				"expression":    updatedExp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
				store.EXPECT().UpdateExpression(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Error - invalid expression ID",
			body: map[string]interface{}{
				"expression_id": "invalid-expression-id",
				"expression":    updatedExp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
				store.EXPECT().UpdateExpression(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Error - empty expression",
			body: map[string]interface{}{
				"expression_id": updatedExp.ExpressionID,
				"expression":    "",
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
				store.EXPECT().UpdateExpression(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Error - database - get expression not found",
			body: map[string]interface{}{
				"expression_id": updatedExp.ExpressionID,
				"expression":    updatedExp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(expstore.Expressions{}, sql.ErrNoRows)
				store.EXPECT().UpdateExpression(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Error - database - get expression internal error",
			body: map[string]interface{}{
				"expression_id": updatedExp.ExpressionID,
				"expression":    updatedExp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(expstore.Expressions{}, sql.ErrConnDone)
				store.EXPECT().UpdateExpression(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Error - expression does not belong to authenticated user",
			body: map[string]interface{}{
				"expression_id": updatedExp.ExpressionID,
				"expression":    updatedExp.Expression,
			},
			bearerToken: anotherAuthValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)
				store.EXPECT().UpdateExpression(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Error - database - update expression not found",
			body: map[string]interface{}{
				"expression_id": updatedExp.ExpressionID,
				"expression":    updatedExp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)

				updateArg := expstore.UpdateExpressionParams{
					ExpressionID: updatedExp.ExpressionID,
					Expression:   updatedExp.Expression,
					UpdatedAt:    updatedExp.UpdatedAt,
				}
				store.EXPECT().
					UpdateExpression(gomock.Any(), EqUpdateExpParams(updateArg)).
					Times(1).Return(expstore.Expressions{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Error - database - update expression internal error",
			body: map[string]interface{}{
				"expression_id": updatedExp.ExpressionID,
				"expression":    updatedExp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)

				updateArg := expstore.UpdateExpressionParams{
					ExpressionID: updatedExp.ExpressionID,
					Expression:   updatedExp.Expression,
					UpdatedAt:    updatedExp.UpdatedAt,
				}
				store.EXPECT().
					UpdateExpression(gomock.Any(), EqUpdateExpParams(updateArg)).
					Times(1).Return(expstore.Expressions{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Error - not authenticated",
			body: map[string]interface{}{
				"expression_id": updatedExp.ExpressionID,
				"expression":    updatedExp.Expression,
			},
			bearerToken: authValue.BearerToken,
			setupAuth:   func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {},
			buildStubs: func(store *mockedexpstore.MockStore) {
				store.EXPECT().GetExpressionByID(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
				store.EXPECT().UpdateExpression(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockedexpstore.NewMockStore(ctrl)
			tc.buildStubs(store)

			config, err := env.NewConfig()
			require.NoError(t, err)

			server, err := expserver.New(config, store, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := "/v1/expressions"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, req, tc.bearerToken)
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestEvaluate(t *testing.T) {
	authValue := authmid.AuthValue{
		BearerToken: authmid.BearerToken1,
		UserID:      "12345",
		Username:    "John Doe",
	}

	exp, qMap := getExpToEvaluate(t, authValue.Username)

	testCases := []struct {
		name          string
		id            string
		queries       map[string]string
		bearerToken   authmid.BearerToken
		setupAuth     func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken)
		buildStubs    func(store *mockedexpstore.MockStore, evaluator *mockedeval.MockEvaluator)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "Happy path",
			id:          exp.ExpressionID.String(),
			queries:     qMap,
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore, evaluator *mockedeval.MockEvaluator) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)
				evaluator.EXPECT().EvalLogicExp(gomock.Any()).Times(1).Return(true)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "Error - invalid expression id",
			id:          "invalid-id",
			queries:     qMap,
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore, evaluator *mockedeval.MockEvaluator) {
				store.EXPECT().GetExpressionByID(gomock.Any(), gomock.Any()).Times(0).Return(expstore.Expressions{}, nil)
				evaluator.EXPECT().EvalLogicExp(gomock.Any()).Times(0).Return(false)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Error - database - get expression not found",
			id:          exp.ExpressionID.String(),
			queries:     qMap,
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore, evaluator *mockedeval.MockEvaluator) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(expstore.Expressions{}, sql.ErrNoRows)
				evaluator.EXPECT().EvalLogicExp(gomock.Any()).Times(0).Return(false)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "Error - database - get expression internal error",
			id:          exp.ExpressionID.String(),
			queries:     qMap,
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore, evaluator *mockedeval.MockEvaluator) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(expstore.Expressions{}, sql.ErrConnDone)
				evaluator.EXPECT().EvalLogicExp(gomock.Any()).Times(0).Return(false)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "Error - missing expression required arguments",
			id:          exp.ExpressionID.String(),
			queries:     nil,
			bearerToken: authValue.BearerToken,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			buildStubs: func(store *mockedexpstore.MockStore, evaluator *mockedeval.MockEvaluator) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(1).Return(exp, nil)
				evaluator.EXPECT().EvalLogicExp(gomock.Any()).Times(0).Return(false)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Error - not authenticated",
			id:          exp.ExpressionID.String(),
			queries:     qMap,
			bearerToken: authValue.BearerToken,
			setupAuth:   func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {},
			buildStubs: func(store *mockedexpstore.MockStore, evaluator *mockedeval.MockEvaluator) {
				store.EXPECT().GetExpressionByID(gomock.Any(), exp.ExpressionID).Times(0).Return(expstore.Expressions{}, nil)
				evaluator.EXPECT().EvalLogicExp(gomock.Any()).Times(0).Return(false)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockedexpstore.NewMockStore(ctrl)
			evaluator := mockedeval.NewMockEvaluator(ctrl)
			tc.buildStubs(store, evaluator)

			config, err := env.NewConfig()
			require.NoError(t, err)

			server, err := expserver.New(config, store, evaluator)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := "/v1/evaluate/" + tc.id

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := req.URL.Query()
			for k, v := range tc.queries {
				if v != "" {
					q.Set(k, v)
				}
			}
			req.URL.RawQuery = q.Encode()

			tc.setupAuth(t, req, tc.bearerToken)
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func getExp(t *testing.T, username string) expstore.Expressions {
	expID, err := uuid.NewRandom()
	require.NoError(t, err)
	require.NotEmpty(t, expID)

	now := time.Now()

	return expstore.Expressions{
		RowID:        1,
		ExpressionID: expID,
		Expression:   "(1 AND 0)",
		Username:     username,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func getExpToEvaluate(t *testing.T, username string) (expstore.Expressions, map[string]string) {
	expID, err := uuid.NewRandom()
	require.NoError(t, err)
	require.NotEmpty(t, expID)

	now := time.Now()
	qMap := map[string]string{
		"x": "1",
		"y": "0",
		"z": "1",
	}

	return expstore.Expressions{
		RowID:        1,
		ExpressionID: expID,
		Expression:   fmt.Sprintf("(%s AND %s) OR %s", "x", "y", "z"),
		Username:     username,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, qMap
}

// addAuthorization adds authorization to the given request
func addAuthorization(
	request *http.Request,
	bearerToken authmid.BearerToken,
	authorizationType string,
) {
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, bearerToken)
	request.Header.Set(authmid.AuthorizationHeaderKey, authorizationHeader)
}

// requireBodyMatchCreate is a helper function to validate the response from the create handler
func requireBodyMatchCreate(t *testing.T, body *bytes.Buffer, exp expstore.Expressions) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var expectedExp, gotExp expmodel.CreateExpressionResponse
	jsonExpectedExp, err := json.Marshal(exp)
	require.NoError(t, err)
	err = json.Unmarshal(jsonExpectedExp, &expectedExp)
	require.NoError(t, err)
	require.NotEmpty(t, expectedExp)

	err = json.Unmarshal(data, &gotExp)
	require.NoError(t, err)
	require.NotEmpty(t, gotExp)

	require.Equal(t, expectedExp.ExpressionID, gotExp.ExpressionID)
	require.Equal(t, expectedExp.Expression, gotExp.Expression)
	require.Equal(t, expectedExp.Username, gotExp.Username)
	require.WithinDuration(t, expectedExp.CreatedAt, gotExp.CreatedAt, time.Second)
	require.WithinDuration(t, expectedExp.UpdatedAt, gotExp.UpdatedAt, time.Second)
}

// requireBodyMatchGet is a helper function to validate the response from the get handler
func requireBodyMatchGet(t *testing.T, body *bytes.Buffer, exp expstore.Expressions) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var expectedExp, gotExp expmodel.CreateExpressionResponse
	jsonExpectedExp, err := json.Marshal(exp)
	require.NoError(t, err)
	err = json.Unmarshal(jsonExpectedExp, &expectedExp)
	require.NoError(t, err)
	require.NotEmpty(t, expectedExp)

	err = json.Unmarshal(data, &gotExp)
	require.NoError(t, err)
	require.NotEmpty(t, gotExp)

	require.Equal(t, expectedExp.ExpressionID, gotExp.ExpressionID)
	require.Equal(t, expectedExp.Expression, gotExp.Expression)
	require.Equal(t, expectedExp.Username, gotExp.Username)
	require.WithinDuration(t, expectedExp.CreatedAt, gotExp.CreatedAt, time.Second)
	require.WithinDuration(t, expectedExp.UpdatedAt, gotExp.UpdatedAt, time.Second)
}

// requireBodyMatchList is a helper function to validate the response from the list handler
func requireBodyMatchList(t *testing.T, body *bytes.Buffer, exps []expstore.Expressions) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var gotExps []expmodel.ListExpressionsResponse
	err = json.Unmarshal(data, &gotExps)
	require.NoError(t, err)

	for i, gotExp := range gotExps {
		require.Empty(t, gotExp.RowID)
		require.Equal(t, exps[i].ExpressionID, gotExp.ExpressionID)
		require.Equal(t, exps[i].Expression, gotExp.Expression)
		require.Equal(t, exps[i].Username, gotExp.Username)
		require.WithinDuration(t, exps[i].CreatedAt, gotExp.CreatedAt, time.Second)
		require.WithinDuration(t, exps[i].UpdatedAt, gotExp.UpdatedAt, time.Second)
	}
}

// requireBodyMatchUpdate is a helper function to validate the response from the update handler
func requireBodyMatchUpdate(t *testing.T, body *bytes.Buffer, exp expstore.Expressions) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var expectedExp, updatedExp expmodel.CreateExpressionResponse
	jsonExpectedExp, err := json.Marshal(exp)
	require.NoError(t, err)
	err = json.Unmarshal(jsonExpectedExp, &expectedExp)
	require.NoError(t, err)
	require.NotEmpty(t, expectedExp)

	err = json.Unmarshal(data, &updatedExp)
	require.NoError(t, err)
	require.NotEmpty(t, updatedExp)

	require.Equal(t, expectedExp.ExpressionID, updatedExp.ExpressionID)
	require.Equal(t, expectedExp.Expression, updatedExp.Expression)
	require.Equal(t, expectedExp.Username, updatedExp.Username)
	require.WithinDuration(t, expectedExp.CreatedAt, updatedExp.CreatedAt, time.Second)
	require.WithinDuration(t, expectedExp.UpdatedAt, updatedExp.UpdatedAt, time.Second)
}
