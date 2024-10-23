package middlewares

import (
	"github.com/gin-gonic/gin"
	"goshop_api/user-web/models"
	"net/http"
)

// IsAdminAuth 检查是否管理员登录
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}

}
