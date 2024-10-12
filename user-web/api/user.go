package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"goshop_api/user-web/global/response"
	"goshop_api/user-web/proto"
	"net/http"
	"time"
)

// HandleGrpcErrorToHttp 将grpc的code转换成http的状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "用户服务不可用",
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误" + e.Message(),
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {
	ip := "127.0.0.1"
	port := 50051

	// 拨号连接用户grpc服务  这里的 Dial 和 WithInsecure 已弃用
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())

	//userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务器失败】",
			"msg", err.Error())
	}

	// 生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)

	// 调用
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    0,
		PSize: 0,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 打印响应数据，用于调试
	//zap.S().Debugw("[GetUserList] 响应数据",
	//	"data", rsp.Data)

	// 返回数据
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       value.Id,
			Nickname: value.NickName,
			//Birthday: time.Time(time.Unix(int64(value.Birthday), 0)).Format("2006-01-02 15:04:05"),
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.Birthday
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile

		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)
}
