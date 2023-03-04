package dao

import (
	"context"
	"fmt"
	"github.com/a754962942/project-project/internal/data/pro"
	"github.com/a754962942/project-project/internal/database/gorms"
)

type ProjectDao struct {
	conn *gorms.GormConn
}

func (p *ProjectDao) FindCollectProjectByMemId(ctx context.Context, memId int64, page int64, size int64, condition string) ([]*pro.ProjectAndMember, int64, error) {
	var pms []*pro.ProjectAndMember
	session := p.conn.Default(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from project where id in (select project_code from project_collection where member_code=?) order by `order` limit ?,?")
	db := session.Raw(sql, memId, index, size)
	db.Scan(&pms)
	var total int64
	query := fmt.Sprintf("member_code=? %s", condition)
	err := session.Model(&pro.ProjectCollection{}).Where(query, memId).Count(&total).Error
	return pms, total, err

}

func (p *ProjectDao) FindProjectByMemId(ctx context.Context, memId int64, page int64, size int64, condition string) ([]*pro.ProjectAndMember, int64, error) {
	var pms []*pro.ProjectAndMember
	session := p.conn.Default(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from project a, project_member b where a.id=b.project_code and b.member_code=?  %s order by `order` limit ?,? ", condition)
	db := session.Raw(sql, memId, index, size)
	db.Scan(&pms)
	var total int64
	query := fmt.Sprintf("select count(*) from project a, project_member b where a.id=b.project_code and b.member_code=? %s", condition)
	tx := session.Raw(query, memId)
	err := tx.Scan(&total).Error
	return pms, total, err
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{
		conn: gorms.New(),
	}
}
