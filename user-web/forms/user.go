package forms

type PasswordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式验证，自定义validator
	Password  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"` // 登录传验证码 长度:5
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}
