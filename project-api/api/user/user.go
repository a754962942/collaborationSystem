package user

import (
	"context"
	common "github.com/a754962942/project-common"
	"github.com/a754962942/project-common/errs"
	loginServiceV1 "github.com/a754962942/project-user/pkg/service/login.service.v1"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HandleUser struct {
}

func New() *HandleUser {
	return &HandleUser{}
}
func (h *HandleUser) getCaptcha(ctx *gin.Context) {
	result := &common.Result{}
	mobile := ctx.PostForm("mobile")
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	captchaResponse, err := LoginServiceClient.GetCaptcha(c, &loginServiceV1.CaptchaMessage{Mobile: mobile})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(captchaResponse.Code))
}
