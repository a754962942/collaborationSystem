package model

import (
	"github.com/a754962942/project-common/errs"
)

var (
	RedisError         = errs.NewError(999, "redis错误")
	DBError            = errs.NewError(998, "DB错误")
	NoLegalMobile      = errs.NewError(10102001, "手机号不合法")
	CaptchaNoExist     = errs.NewError(10102002, "验证码不存在或已过期")
	CaptchaError       = errs.NewError(10102002, "验证码不正确")
	EmailExist         = errs.NewError(10102003, "邮箱已注册")
	NameExist          = errs.NewError(10102004, "账号已注册")
	MobileExist        = errs.NewError(10102005, "手机号已注册")
	AccountAndPwdEroor = errs.NewError(10102006, "账号或密码不正确")
	NoLogin            = errs.NewError(10102007, "未登录")
)
