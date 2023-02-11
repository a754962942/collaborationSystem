package model

import (
	"github.com/a754962942/project-common/errs"
)

var (
	NoLegalMobile = errs.NewError(2001, "手机号不合法")
)
