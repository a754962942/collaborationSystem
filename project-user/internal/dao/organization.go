package dao

import (
	"context"
	"github.com/a754962942/project-user/internal/data/organization"
	"github.com/a754962942/project-user/internal/database"
	"github.com/a754962942/project-user/internal/database/gorms"
)

type Organization struct {
	conn *gorms.GormConn
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
