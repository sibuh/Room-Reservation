package auth

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type Authorizer interface {
	CreateRole(c *gin.Context)
	CreatePermission(c *gin.Context)
}

type authorizer struct {
	logger *slog.Logger
	e      casbin.IEnforcer
}

func NewAuthorizer(logger *slog.Logger, e casbin.IEnforcer) Authorizer {
	return &authorizer{
		logger: logger,
		e:      e,
	}
}

func (a *authorizer) CreateRole(c *gin.Context) {
	
}

func (a *authorizer) CreatePermission(c *gin.Context) {

}
