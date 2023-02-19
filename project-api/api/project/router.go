package project

import (
	"github.com/a754962942/project-api/router"
	"github.com/gin-gonic/gin"
)

type RouterPeoject struct {
}

func init() {
	ru := &RouterPeoject{}
	router.Register(ru)
}
func (*RouterPeoject) Route(r *gin.Engine) {
	//初始化grpc的客户端连接
	InitRpcUserClient()
	project := New()
	r.POST("/project/index", project.index)

}
