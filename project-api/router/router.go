package router

import (
	"fmt"
	"github.com/a754962942/project-common/logs"
	"github.com/a754962942/project-user/config"
	loginServiceV1 "github.com/a754962942/project-user/pkg/service/login.service.v1"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net"
)

var routers []Router

func Register(ro ...Router) {
	routers = append(routers, ro...)

}

type Router interface {
	Route(r *gin.Engine)
}
type RegisterRouter struct {
}

func New() *RegisterRouter {
	return &RegisterRouter{}
}
func (*RegisterRouter) Route(ro Router, r *gin.Engine) {
	ro.Route(r)
}
func InitRouter(r *gin.Engine) {
	//1.手动注册路由
	//新增模块需要在此处进行挂载
	//rg := New()
	//rg.Route(&user.RouterUser{}, r)
	//2.自动注册路由
	//新增模块需要调用register对路由进行append
	//各模块需在init中进行append操作。
	for _, router := range routers {
		router.Route(r)
	}
}

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(server *grpc.Server)
}

func RegisterGrpc() *grpc.Server {
	g := &gRPCConfig{
		Addr: config.C.Gc.Addr,
		RegisterFunc: func(g *grpc.Server) {
			loginServiceV1.RegisterLoginServiceServer(g, loginServiceV1.New())
		},
	}
	s := grpc.NewServer()
	g.RegisterFunc(s)
	listen, err := net.Listen("tcp", g.Addr)
	if err != nil {
		logs.LG.Info("connot listen")
	}
	go func() {
		err := s.Serve(listen)
		if err != nil {
			logs.LG.Error(fmt.Sprintf("Server started error %s", err))
			return
		}
	}()
	return s
}
