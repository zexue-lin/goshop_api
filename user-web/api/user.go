package api

import (
	"context"
	"fmt"
	"goshop_api/user-web/forms"
	"goshop_api/user-web/global"
	"goshop_api/user-web/global/response"
	"goshop_api/user-web/middlewares"
	"goshop_api/user-web/models"
	"goshop_api/user-web/proto"
	"net/http"

	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// 去除错误信息中的表单名称
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err // field -> 要查找的字符串， . => 查找点号的位置
	}
	return rsp
}

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

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)), // Translate 本质上返回的就是一个map[string]string
	})
	return
}

func GetUserList(ctx *gin.Context) {

	// 拨号连接用户grpc服务  这里的 Dial 和 WithInsecure 已弃用
	// userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())

	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务器失败】",
			"msg", err.Error())
	}

	// 生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)

	// 调用, 使用了ctx上下文对象来获取前端传来的查询参数，proto中定义的是uint32类型
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
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
		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)
}

// 登录的逻辑
func PasswordLogin(c *gin.Context) {
	// 表单验证
	passwordLoginForm := forms.PasswordLoginForm{}

	// 下面的if逻辑再每次做表单验证的时候都会使用到，比较长，可以封装起来
	// 刚才已经做过类似的操作 就是函数 HandleGrpcErrorToHttp
	// 把注释去掉就变得很简单了
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		// // 把 err 类型转换一下
		// errs, ok := err.(validator.ValidationErrors)
		// if !ok {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"message": err.Error(),
		// 	})
		// }
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"error": removeTopStruct(errs.Translate(global.Trans)), // Translate 本质上返回的就是一个map[string]string
		// })
		HandleValidatorError(c, err)
		return
	}

	// 拨号连接用户grpc服务  这里的 Dial 和 WithInsecure 已弃用
	// userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())

	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务器失败】",
			"msg", err.Error())
	}

	// 生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)

	// 登录的逻辑，前面是判断，检验
	if rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		// 只是查询到了用户而已，并没有检查密码,
		if passRsp, pasErr := userSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: rsp.Password,
		}); pasErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"msg": "登录失败",
			})
		} else {
			if passRsp.Success {
				// 生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               // 前面的生成时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer:    "lattiex",                       // 签名机构
					},
				}

				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return 
				}

				c.JSON(http.StatusOK, gin.H{
					"id": rsp.Id,
					"nick_name": rsp.NickName,
					"token": token,
					"expired_at": (time.Now().Unix() + 60*60*24*30)*1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登陆失败",
				})
			}
		}
	}
}
