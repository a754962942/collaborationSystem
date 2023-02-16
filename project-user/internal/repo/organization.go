package repo

import (
	"context"
	"github.com/a754962942/project-user/internal/data/organization"
	"github.com/a754962942/project-user/internal/database"
)

type Organization interface {
	SaveOrganization(conn database.DbConn, ctx context.Context, organization *organization.Organization) error
	FindOrganizationByMemId(ctx context.Context, memId int64) ([]*organization.Organization, error)
}
