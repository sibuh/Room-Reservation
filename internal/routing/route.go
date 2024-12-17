package routing

import "github.com/gin-gonic/gin"

var RouteDetails = []struct {
	Path        string
	Method      string
	Middlewares []gin.HandlerFunc
	handler     gin.HandlerFunc
}{
	{},
}

func RegisterRoutes() {
	var r = gin.Default()

	r.POST("/", gin.HandlerFunc(func(ctx *gin.Context) {}))
}
