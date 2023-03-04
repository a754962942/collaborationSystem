package project_service_v1

import (
	"context"
	"github.com/a754962942/project-common/encrypts"
	"github.com/a754962942/project-common/errs"
	"github.com/a754962942/project-common/tms"
	"github.com/a754962942/project-grpc/project"
	"github.com/a754962942/project-project/internal/dao"
	"github.com/a754962942/project-project/internal/data/menu"
	"github.com/a754962942/project-project/internal/data/pro"
	"github.com/a754962942/project-project/internal/database/tran"
	"github.com/a754962942/project-project/internal/repo"
	"github.com/a754962942/project-project/pkg/model"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type ProjectService struct {
	project.UnimplementedRegisterServiceServer
	cache               repo.Cache
	transaction         tran.Transaction
	menuRepo            repo.MenuRepo
	projectRepo         repo.ProjectRepo
	projectTemplateRepo repo.ProjectTemplateRepo
}

func New() *ProjectService {
	return &ProjectService{
		cache:               dao.Rc,
		transaction:         dao.NewTransactionImpl(),
		menuRepo:            dao.NewMenuDao(),
		projectRepo:         dao.NewProjectDao(),
		projectTemplateRepo: dao.NewProjectTemplateDao(),
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
	var pms []*pro.ProjectAndMember
	var total int64
	var err error
	if msg.SelectBy == "" || msg.SelectBy == "my" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, id, page, size, "")
	}
	if msg.SelectBy == "collect" {
		pms, total, err = p.projectRepo.FindCollectProjectByMemId(ctx, id, page, size, "a")
	}
	if msg.SelectBy == "archive" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, id, page, size, "and archive=1 ")
	}
	if msg.SelectBy == "deleted" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, id, page, size, "and deleted=1 ")
	}
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
		pam := pro.ToMap(pms)[v.Id]
		v.AccessControlType = pam.GetAccessControlType()
		v.OrganizationCode, _ = encrypts.EncryptInt64(pam.OrganizationCode, model.AESKey)
		v.JoinTime = tms.FormatByMill(pam.JoinTime)
		v.OwnerName = msg.MemberName
		v.Order = int32(pam.Order)
		v.CreateTime = tms.FormatByMill(pam.CreateTime)
	}
	return &project.MyProjectResponse{Pm: messages, Total: total}, nil
}
func (ps *ProjectService) FindProjectTemplate(ctx context.Context, msg *project.ProjectRpcMessage) (*project.ProjectTemplateResponse, error) {
	//1.根据viewType查询项目模板表 得到list
	code := msg.OrganizationCode

	var pts []pro.ProjectTemplate
	var total int64
	var err error
	if msg.ViewType == -1 {
		ps.projectTemplateRepo.FindProjectTemplateAll(ctx)
	}
	if msg.ViewType == 0 {

	}
	if msg.ViewType == 1 {

	}
	//2.模型转换，拿到模板id列表后去任务步骤模板表进行查询
	//3.组装数据
}
