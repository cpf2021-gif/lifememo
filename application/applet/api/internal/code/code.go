package code

import "lifememo/pkg/xcode"

var (
	EmailEmpty         = xcode.New(10000, "邮箱不能为空")
	EmailFormatError   = xcode.New(10001, "邮箱格式错误")
	RegisterEmailExist = xcode.New(10002, "注册邮箱已存在")

	RegisterNameEmpty = xcode.New(10003, "注册用户名不能为空")
	RegisterNameLimit = xcode.New(10004, "注册用户名长度不能超过20个字符")

	VerificationCodeEmpty   = xcode.New(10005, "验证码不能为空")
	VerificationCodeError   = xcode.New(10006, "验证码错误")
	VerificationCodeExpired = xcode.New(10007, "验证码已过期")
	VerificationCodeLimit   = xcode.New(10008, "验证码发送次数超过限制")
	SendEmailError          = xcode.New(10009, "发送验证码失败, 请稍后重试")
)
