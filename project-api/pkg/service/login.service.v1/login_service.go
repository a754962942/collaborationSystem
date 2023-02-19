package login_service_v1

import (
	"context"
	"errors"
	"github.com/a754962942/project-api/pkg/dao"
	"github.com/a754962942/project-api/pkg/repo"
	common "github.com/a754962942/project-common"
	"github.com/a754962942/project-common/logs"

	"log"
	"time"
)

type LoginService struct {
	UnimplementedLoginServiceServer
	cache repo.Cache
}

func New() *LoginService {
	return &LoginService{
		cache: dao.Rc,
	}
}
func (ls *LoginService) GetCaptcha(ctx context.Context, msg *CaptchaMessage) (*CaptchaResponse, error) {

	//1. 获取参数
	mobile := msg.Mobile
	//2. 校验参数
	if !common.VerifyMobile(mobile) {
		return nil, errors.New("手机号不合法")
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
		err := ls.cache.Put(c, "REGISTER_"+mobile, code, 15*time.Minute)
		if err != nil {
			log.Printf("验证码存入redis出错,cause by:%s\n", err)
			return
		}
		log.Printf("将手机号和验证码存redis成功：REGISTER_%s:%s\n", mobile, code)
	}()
	return &CaptchaResponse{}, nil
}
