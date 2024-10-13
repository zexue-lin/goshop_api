package forms

type PasswordLoginForm struct{
	Mobile string `form:"mobile" json:"mobile" binding:"required"`  // 手机号码格式验证，自定义validator
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}