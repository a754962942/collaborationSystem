package dao

import (
	"context"
	"github.com/a754962942/project-project/internal/data/pro"
	"github.com/a754962942/project-project/internal/database/gorms"
)

type ProjectTemplateDao struct {
	conn *gorms.GormConn
}

func (p ProjectTemplateDao) FindProjectTemplateSystem(ctx context.Context, page int64, size int64) ([]pro.ProjectTemplate, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectTemplateDao) FindProjectTemplateCustom(ctx context.Context, memId int64, organizationCode int64, page int64, size int64) ([]pro.ProjectTemplate, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectTemplateDao) FindProjectTemplateAll(ctx context.Context, organizationCode int64, page int64, size int64) ([]pro.ProjectTemplate, int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewProjectTemplateDao() *ProjectTemplateDao {
	return &ProjectTemplateDao{
		conn: gorms.New(),
	}
}
