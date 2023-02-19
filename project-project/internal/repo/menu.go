package repo

import (
	"context"
	"github.com/a754962942/project-project/internal/data/menu"
)

type MenuRepo interface {
	FindMenus(ctx context.Context) ([]*menu.ProjectMenu, error)
}
