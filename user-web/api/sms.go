package api

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"
	"goshop_api/user-web/forms"
	"goshop_api/user-web/global"
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
	// 表单验证 - 发送手机验证码
	SendSmsForm := forms.SendSmsForm{}

	if err := ctx.ShouldBind(&SendSmsForm); err != nil {

		HandleValidatorError(ctx, err)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecret)
	if err != nil {
		panic(err)
	}

	smsCode := GenerateSmsCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.Version = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = SendSmsForm.Mobile
	request.QueryParams["SignName"] = "Go商城"
	request.QueryParams["TemplateCode"] = "SMS_181850725"
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}"
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	if err != nil {
		fmt.Println(err.Error())
	}

	// 将验证码保存起来 - 放在redis里面
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%d:%s", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	rdb.Set(context.Background(), SendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.AliSmsInfo.Expire)*time.Second)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
