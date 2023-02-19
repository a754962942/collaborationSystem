package project

import (
	"context"
	common "github.com/a754962942/project-common"
	"github.com/a754962942/project-common/errs"
	"github.com/a754962942/project-grpc/project"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HandleProject struct {
}

func (p *HandleProject) index(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.IndexMessage{}
	indexResponse, err := RegisterServiceClient.Index(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(indexResponse.Menus))
}

func New() *HandleProject {
	return &HandleProject{}
}
