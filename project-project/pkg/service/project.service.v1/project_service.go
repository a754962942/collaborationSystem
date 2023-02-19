package project_service_v1

import (
	"context"
	"github.com/a754962942/project-common/errs"
	"github.com/a754962942/project-grpc/project"
	"github.com/a754962942/project-project/internal/dao"
	"github.com/a754962942/project-project/internal/data/menu"
	"github.com/a754962942/project-project/internal/database/tran"
	"github.com/a754962942/project-project/internal/repo"
	"github.com/a754962942/project-project/pkg/model"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type ProjectService struct {
	project.UnimplementedRegisterServiceServer
	cache       repo.Cache
	transaction tran.Transaction
	menuRepo    repo.MenuRepo
}

func New() *ProjectService {
	return &ProjectService{
		cache:       dao.Rc,
		transaction: dao.NewTransactionImpl(),
		menuRepo:    dao.NewMenuDao(),
	}
}
func (p *ProjectService) Index(context.Context, *project.IndexMessage) (*project.IndexResponse, error) {
	pms, err := p.menuRepo.FindMenus(context.Background())
	if err != nil {
		zap.L().Error("Index DB FindMenus error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	childs := menu.CovertChild(pms)
	var mms []*project.MenuMessage
	_ = copier.Copy(&mms, childs)
	return &project.IndexResponse{Menus: mms}, nil
}
