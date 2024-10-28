package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式验证，自定义validator
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`
	// 1.注册发送短信验证码 和 动态验证码登录发送验证码
	// 1代表注册 2代表动态验证码登录
}
