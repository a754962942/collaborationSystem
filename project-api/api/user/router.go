package user

import (
	"github.com/a754962942/project-user/router"
	"github.com/gin-gonic/gin"
)

type RouterUser struct {
}

func init() {
	ru := &RouterUser{}
	router.Register(ru)
}
func (*RouterUser) Route(r *gin.Engine) {
	//初始化grpc的客户端连接
	InitRpcUserClient()
	user := New()
	r.POST("/project/login/getCaptcha", user.getCaptcha)
	r.POST("/project/login/register", user.register)
	r.POST("/project/login", user.login)
}
