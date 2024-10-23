package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"net/http"
	"strings"
	"time"
)

func GenerateSmsCode(width int) string {
	// 生成长度为width的短信验证码
	numernic := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numernic)
	rand.Seed(uint64(time.Now().UnixNano()))

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numernic[rand.Intn(r)])
	}
	return sb.String()

}
func SendSms(ctx *gin.Context) {

	zap.S().Errorf("生成验证码错误:%s", err.Error())
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"msg": "生成验证码错误",
	})

	// 将验证码保存起来 - 放在redis里面

}
