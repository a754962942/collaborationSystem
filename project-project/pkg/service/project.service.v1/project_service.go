package project_service_v1

import (
	"context"
	"github.com/a754962942/project-common/encrypts"
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
	projectRepo repo.ProjectRepo
}

func New() *ProjectService {
	return &ProjectService{
		cache:       dao.Rc,
		transaction: dao.NewTransactionImpl(),
		menuRepo:    dao.NewMenuDao(),
		projectRepo: dao.NewProjectDao(),
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
func (p *ProjectService) FindProjectByMemId(ctx context.Context, msg *project.ProjectRpcMessage) (*project.MyProjectResponse, error) {
	id := msg.MemberId
	page := msg.Page
	size := msg.PageSize
	pms, total, err := p.projectRepo.FindProjectByMemId(ctx, id, page, size)
	if err != nil {
		zap.L().Error("Project FindProjectByMemberId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if pms == nil {
		return &project.MyProjectResponse{Pm: []*project.ProjectMessage{}, Total: total}, nil
	}
	messages := []*project.ProjectMessage{}
	_ = copier.Copy(&messages, pms)
	for _, v := range messages {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)
	}
	return &project.MyProjectResponse{Pm: messages, Total: total}, nil
}
