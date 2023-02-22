package dao

import (
	"context"
	"github.com/a754962942/project-project/internal/data/pro"
	"github.com/a754962942/project-project/internal/database/gorms"
)

type ProjectDao struct {
	conn *gorms.GormConn
}

func (p *ProjectDao) FindProjectByMemId(ctx context.Context, memId int64, page int64, size int64) ([]*pro.ProjectAndMember, int64, error) {
	var pms []*pro.ProjectAndMember
	session := p.conn.Default(ctx)
	index := (page - 1) * size
	db := session.Raw("select * from project a,project_member b where a.id = b.project_code and b.member_code =? limit ?,?", memId, index, size)
	db.Scan(&pms)
	var total int64
	err := session.Model(&pro.ProjectMember{}).Where("member_code =?", memId).Count(&total).Error
	return pms, total, err
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{
		conn: gorms.New(),
	}
}
