package middlewares

import (
	"errors"
	"net/http"

	errApi "chat-system/cmd/api/errors"
	"chat-system/cmd/api/security"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, errApi.NewErrAPIUnauthorized(errors.New("authorization header is missing")))
			c.Abort()
			return
		}

		claims, err := security.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errApi.NewErrAPIUnauthorized(errors.New("invalid or expired token")))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}
