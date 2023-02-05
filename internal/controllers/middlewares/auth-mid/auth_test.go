package authmid_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	authmid "github.com/gmaschi/log-exp-eval/internal/controllers/middlewares/auth-mid"
	expserver "github.com/gmaschi/log-exp-eval/internal/servers/expressions"
	"github.com/gmaschi/log-exp-eval/pkg/tools/config/env"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		bearerToken   authmid.BearerToken
		setupAuth     func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "Happy path",
			bearerToken: authmid.BearerToken1,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "Error - No Authorization",
			bearerToken: authmid.BearerToken1,
			setupAuth:   func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:        "Error - Unsupported Authorization Type",
			bearerToken: authmid.BearerToken1,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, "unsupported")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:        "Error - Invalid Authorization Format",
			bearerToken: authmid.BearerToken1,
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, "")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:        "Error - Invalid Bearer Token",
			bearerToken: authmid.BearerToken("invalid-token"),
			setupAuth: func(t *testing.T, request *http.Request, bearerToken authmid.BearerToken) {
				addAuthorization(request, bearerToken, authmid.AuthorizationTypeBearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config, err := env.NewConfig()
			require.NoError(t, err)
			server, err := expserver.New(config, nil, nil)
			require.NoError(t, err)

			authPath := "/exp"

			server.Router.GET(
				authPath,
				authmid.AuthMiddleware(),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, map[string]interface{}{})
				},
			)

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, tc.bearerToken)
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
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
