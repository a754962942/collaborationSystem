package dao

import (
	"context"
	"github.com/a754962942/project-project/internal/data/menu"
	"github.com/a754962942/project-project/internal/database/gorms"
)

type MenuDao struct {
	conn *gorms.GormConn
}

func (m *MenuDao) FindMenus(ctx context.Context) (pms []*menu.ProjectMenu, err error) {
	session := m.conn.Default(ctx)
	err = session.Find(&pms).Error
	return pms, err
}

func NewMenuDao() *MenuDao {
	return &MenuDao{conn: gorms.New()}
}
