package project

import (
	"github.com/a754962942/project-api/api/midd"
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
	InitRpcRegisterClient()
	project := New()
	r.Use(midd.TokenVerify())
	group := r.Group("/project/index")
	group.POST("", project.index)
	group1 := r.Group("/project/project")
	group1.POST("/selfList", project.myProjectList)
	group1.POST("", project.myProjectList)
}
