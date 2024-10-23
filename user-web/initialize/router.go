package initialize

import (
	"goshop_api/user-web/middlewares"
	"goshop_api/user-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	// 配置跨域
	Router.Use(middlewares.CORS())

	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return Router
}
