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
	"github.com/a754962942/project-project/internal/data/task"
	"github.com/a754962942/project-project/internal/database/tran"
	"github.com/a754962942/project-project/internal/repo"
	"github.com/a754962942/project-project/pkg/model"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"strconv"
)

type ProjectService struct {
	project.UnimplementedRegisterServiceServer
	cache                  repo.Cache
	transaction            tran.Transaction
	menuRepo               repo.MenuRepo
	projectRepo            repo.ProjectRepo
	projectTemplateRepo    repo.ProjectTemplateRepo
	taskStagesTemplateRepo repo.TaskStagesTemplateRepo
}

func New() *ProjectService {
	return &ProjectService{
		cache:                  dao.Rc,
		transaction:            dao.NewTransactionImpl(),
		menuRepo:               dao.NewMenuDao(),
		projectRepo:            dao.NewProjectDao(),
		projectTemplateRepo:    dao.NewProjectTemplateDao(),
		taskStagesTemplateRepo: dao.NewTaskStagesTemplateDao(),
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
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	page := msg.Page
	pageSize := msg.PageSize
	var pts []pro.ProjectTemplate
	var total int64
	var err error
	if msg.ViewType == -1 {
		pts, total, err = ps.projectTemplateRepo.FindProjectTemplateAll(ctx, organizationCode, page, pageSize)
		if err != nil {
			zap.L().Error("project FindProjectTemplate FindProjectTemplateAll error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	if msg.ViewType == 0 {
		pts, total, err = ps.projectTemplateRepo.FindProjectTemplateCustom(ctx, msg.MemberId, organizationCode, page, pageSize)
		if err != nil {
			zap.L().Error("project FindProjectTemplate FindProjectTemplateCustom error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	if msg.ViewType == 1 {
		pts, total, err = ps.projectTemplateRepo.FindProjectTemplateSystem(ctx, page, pageSize)
		if err != nil {
			zap.L().Error("project FindProjectTemplate FindProjectTemplateSystem error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}

	//2.模型转换，拿到模板id列表后去任务步骤模板表进行查询
	tsts, err := ps.taskStagesTemplateRepo.FindInProTemIds(ctx, pro.ToProjectTemplateIds(pts))
	if err != nil {
		zap.L().Error("project FindProjectTemplate FindInProTemIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var ptas []*pro.ProjectTemplateAll
	for _, v := range pts {
		//该谁做的一定要交出去
		ptas = append(ptas, v.Convert(task.CovertProjectMap(tsts)[v.Id]))
	}
	//3.组装数据
	var pmMsgs []*project.ProjectTemplateMessage
	copier.Copy(pmMsgs, ptas)

	return &project.ProjectTemplateResponse{Ptm: pmMsgs, Total: total}, nil
}
