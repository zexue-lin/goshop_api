package router

import (
	"goshop_api/user-web/api"
	"goshop_api/user-web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 注册用户相关的路由
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PasswordLogin)
	}

}
