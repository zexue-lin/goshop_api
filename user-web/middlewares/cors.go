package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token, x-token")
		c.Header("ACCESS-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT,PATCH")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin, Content-Type, Access-Control-Allow-Headers")
		c.Header("ACCESS-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
	/*
		204 No Content表示请求成功但不返回任何内容。这在处理预检请求时非常合适，
		因为OPTIONS请求的目的只是确认服务器是否允许跨域，不需要返回实际的数据。

		预检请求的目标仅仅是为了验证CORS策略，而不是获取实际的数据!!!
	*/
}
