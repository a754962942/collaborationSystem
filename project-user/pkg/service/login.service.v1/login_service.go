package login_service_v1

import (
	"context"
	common "github.com/a754962942/project-common"
	"github.com/a754962942/project-common/encrypts"
	"github.com/a754962942/project-common/errs"
	"github.com/a754962942/project-common/logs"
	"github.com/a754962942/project-grpc/user/login"
	"github.com/a754962942/project-user/internal/dao"
	"github.com/a754962942/project-user/internal/data/member"
	"github.com/a754962942/project-user/internal/data/organization"
	"github.com/a754962942/project-user/internal/database"
	"github.com/a754962942/project-user/internal/database/tran"
	"github.com/a754962942/project-user/internal/repo"
	"github.com/a754962942/project-user/pkg/model"
	"go.uber.org/zap"
	"log"
	"time"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache            repo.Cache
	memberRepo       repo.MemberRepo
	organizationRepo repo.Organization
	transaction      tran.Transaction
}

func New() *LoginService {
	return &LoginService{
		cache:            dao.Rc,
		memberRepo:       dao.NerMemberDao(),
		organizationRepo: dao.NerOrganization(),
		transaction:      dao.NerTransactionImpl(),
	}
}
func (ls *LoginService) GetCaptcha(ctx context.Context, msg *login.CaptchaMessage) (*login.CaptchaResponse, error) {
	//1. 获取参数
	mobile := msg.Mobile
	//2. 校验参数
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	//3. 生成验证码(随机4或6位)
	code := "123456"
	//4. 调用短信平台(第三方 放入go协程中执行，接口可以快速响应)
	go func() {
		time.Sleep(1 * time.Second)
		logs.LG.Info("短信平台调用成功,发送短信 INFO")
		logs.LG.Debug("短信平台调用成功,发送短信 debug")
		logs.LG.Error("短信平台调用成功,发送短信 error")
		//redis 假设后续缓存可能存在mysql当中，也可能存在mongo中，也可能存在memcache当中
		//5. 存储验证码 redis 过期时间15min
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := ls.cache.Put(c, model.RegisterRedisKey+mobile, code, 15*time.Minute)
		if err != nil {
			log.Printf("验证码存入redis出错,cause by:%s\n", err)
			return
		}
		log.Printf("将手机号和验证码存redis成功：REGISTER_%s:%s\n", mobile, code)
	}()
	return &login.CaptchaResponse{Code: code}, nil
}
func (ls *LoginService) Register(ctx context.Context, msg *login.RegisterMessage) (*login.RegisterResponse, error) {
	c := context.Background()
	//	1.校验参数
	//	2.校验验证码
	redisCode, err := ls.cache.Get(c, model.RegisterRedisKey+msg.Mobile)
	if err != nil {
		zap.L().Error("Register redis get error", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != msg.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	//	3.校验业务逻辑(邮箱是否被注册|账号是否被注册|手机号是否被注册)
	exist, err := ls.memberRepo.GetMemberByEmail(c, msg.Email)
	if err != nil {
		zap.L().Error("Register DB get error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist)
	}
	exist, err = ls.memberRepo.GetMemberByAccount(c, msg.Name)
	if err != nil {
		zap.L().Error("Register DB get error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.NameExist)
	}
	exist, err = ls.memberRepo.GetMemberByMobile(c, msg.Mobile)
	if err != nil {
		zap.L().Error("Register DB get error", zap.Error(err))
		return nil, errs.GrpcError(model.CaptchaNoExist)
	}
	if exist {
		return nil, errs.GrpcError(model.MobileExist)
	}
	//	4.执行业务 将数据存入member表 生成一个数据存入organization表
	pwd := encrypts.Md5(msg.Password)
	mem := &member.Member{
		Account:       msg.Name,
		Password:      pwd,
		Name:          msg.Name,
		Mobile:        msg.Mobile,
		Email:         msg.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        model.Normal,
	}
	err = ls.transaction.Action(func(conn database.DbConn) error {
		err = ls.memberRepo.SaveMember(conn, c, mem)
		if err != nil {
			zap.L().Error("SavaMember get error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		//存入组织
		org := &organization.Organization{
			Name:       mem.Name + "个人项目",
			MemberId:   mem.Id,
			CreateTime: time.Now().UnixMilli(),
			Personal:   model.Personal,
			Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		}
		err = ls.organizationRepo.SaveOrganization(conn, c, org)
		if err != nil {
			zap.L().Error("register SaveOrganization db err", zap.Error(err))
			return model.DBError
		}
		return nil
	})
	if err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	//	5.返回
	return &login.RegisterResponse{}, nil
}
