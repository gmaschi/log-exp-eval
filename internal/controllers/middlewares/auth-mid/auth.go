package authmid

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gmaschi/log-exp-eval/pkg/tools/parse"
)

type BearerToken string

func (bt BearerToken) String() string {
	return string(bt)
}

// AuthValue defines the authorization type containing all the information related to each of the roles
type AuthValue struct {
	BearerToken BearerToken
	UserID      string
	Username    string
}

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"

	BearerToken1 BearerToken = "74edf612f393b4eb01fbc2c29dd96671"
	BearerToken2 BearerToken = "d88b4b1e77c70ba780b56032db1c259b"
)

var (
	authMap = map[BearerToken]AuthValue{
		BearerToken1: {
			BearerToken: BearerToken1,
			UserID:      "12345",
			Username:    "John Doe",
		},
		BearerToken2: {
			BearerToken: BearerToken2,
			UserID:      "98765",
			Username:    "Jane Doe",
		},
	}
)

// AuthMiddleware defines the authentication middleware that it's going to validate the request
// before forwarding it to the handlers
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, parse.ErrorAsJSON(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, parse.ErrorAsJSON(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization format %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, parse.ErrorAsJSON(err))
			return
		}

		accessToken := fields[1]
		authValue, ok := authMap[BearerToken(accessToken)]
		if !ok {
			err := fmt.Errorf("invalid bearer token provided %s", accessToken)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, parse.ErrorAsJSON(err))
			return
		}

		ctx.Set(AuthorizationPayloadKey, authValue)
		ctx.Next()
	}
}
