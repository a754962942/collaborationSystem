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
	group := r.Group("/project")
	group.Use(midd.TokenVerify())
	group.POST("/index", project.index)
	group.POST("/project/selfList", project.myProjectList)
	group.POST("/project", project.myProjectList)
	group.POST("/project_template", project.projectTemplate)

}
