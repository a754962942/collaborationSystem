package user

import (
	"context"
	"github.com/a754962942/project-api/pkg/model/user"
	common "github.com/a754962942/project-common"
	"github.com/a754962942/project-common/errs"
	"github.com/a754962942/project-grpc/user/login"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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
	captchaResponse, err := LoginServiceClient.GetCaptcha(c, &login.CaptchaMessage{Mobile: mobile})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(captchaResponse.Code))
}
func (h *HandleUser) register(c *gin.Context) {
	//1.接受参数 参数模型
	result := &common.Result{}
	req := user.RegisterReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式错误"))
		return
	}
	//2.校验参数 判断参数是否合法
	if err := req.Verify(); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, err.Error()))
		return
	}
	//3.调用user grpc服务 获取响应
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	msg := &login.RegisterMessage{}
	err := copier.Copy(msg, req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, err.Error()))
	}
	_, err = LoginServiceClient.Register(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(""))

}
