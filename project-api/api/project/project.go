package project

import (
	"context"
	"fmt"
	"github.com/a754962942/project-api/pkg/model"
	"github.com/a754962942/project-api/pkg/model/menu"
	"github.com/a754962942/project-api/pkg/model/pro"
	common "github.com/a754962942/project-common"
	"github.com/a754962942/project-common/errs"
	"github.com/a754962942/project-grpc/project"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
	"time"
)

type HandleProject struct {
}

func (p *HandleProject) index(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.IndexMessage{}
	indexResponse, err := ProjectServiceClient.Index(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	var mn []*menu.Menu
	err = copier.Copy(&mn, indexResponse.Menus)
	fmt.Println(err)
	c.JSON(http.StatusOK, result.Success(mn))
}

func (p *HandleProject) myProjectList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberstr, _ := c.Get("memberId")
	memName, _ := c.Get("memberName")
	memberId := memberstr.(int64)
	page := &model.Page{}
	page.Bind(c)
	selectBy := c.PostForm("selectBy")
	msg := &project.ProjectRpcMessage{MemberId: memberId, MemberName: memName.(string), PageSize: page.PageSize, Page: page.Page, SelectBy: selectBy}
	response, err := ProjectServiceClient.FindProjectByMemId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	var pms []*pro.ProjectAndMember
	copier.Copy(&pms, response.Pm)
	if pms == nil {
		pms = []*pro.ProjectAndMember{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms,
		"total": response.Total,
	}))
}

func (p *HandleProject) projectTemplate(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	organizationCode := c.GetString("organizationCode")
	var page = &model.Page{}
	page.Bind(c)
	viewTypeStr := c.PostForm("viewType")
	viewType, _ := strconv.ParseInt(viewTypeStr, 10, 64)
	projectTemplateRsp, err := ProjectServiceClient.FindProjectTemplate(ctx,
		&project.ProjectRpcMessage{
			MemberId:         memberId,
			MemberName:       memberName,
			OrganizationCode: organizationCode,
			Page:             page.Page,
			PageSize:         page.PageSize,
			ViewType:         int32(viewType)})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var pts []*pro.ProjectTemplate
	copier.Copy(&pts, projectTemplateRsp.Ptm)
	if pts == nil {
		pts = []*pro.ProjectTemplate{}
	}
	c.JSON(http.StatusOK, result.Success(pts))
}

func New() *HandleProject {
	return &HandleProject{}
}
