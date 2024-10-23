package router

import (
	"goshop_api/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("chptcha", api.GetCaptcha)
		BaseRouter.POST("send_sms", api.GetCaptcha)
	}
}
