package middleware

import (
	"errors"
	"net/http"
	"reservation/pkg/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const Bearer = "Bearer"

var (
	ErrAuthorizedAccess = errors.New("")
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn := c.Request.Header.Get("Authorization")
		if tkn == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("no token not provided"))
		}
		slicedToken := strings.Split(tkn, " ")
		if slicedToken[0] != Bearer {
			c.AbortWithError(http.StatusUnauthorized, errors.New("token is not of type bearer"))
		}
		payload := token.VerifyToken(slicedToken[1])

	}
}
