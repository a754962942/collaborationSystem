package model

import (
	"github.com/a754962942/project-common/errs"
)

var (
	RedisError    = errs.NewError(-100, "redis错误")
	NoLegalMobile = errs.NewError(10102001, "手机号不合法")
	CaptchaError  = errs.NewError(10102002, "验证码不正确")
)
