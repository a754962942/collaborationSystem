package dao

import (
	"context"
	"github.com/a754962942/project-user/internal/data/organization"
	"github.com/a754962942/project-user/internal/database"
	"github.com/a754962942/project-user/internal/database/gorms"
	"gorm.io/gorm"
)

type Organization struct {
	conn *gorms.GormConn
}

func (o *Organization) FindOrganizationByMemId(ctx context.Context, memId int64) ([]*organization.Organization, error) {
	orgs := make([]*organization.Organization, 0)
	err := o.conn.Default(ctx).Where("member_id=?", memId).Find(&orgs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return orgs, err
}

func NerOrganization() *Organization {
	return &Organization{
		conn: gorms.New(),
	}
}
func (m *Organization) SaveOrganization(conn database.DbConn, ctx context.Context, organization *organization.Organization) error {
	m.conn = conn.(*gorms.GormConn)
	return m.conn.Tx(ctx).Create(organization).Error
}
