package router

import (
	"fmt"
	"github.com/a754962942/project-common/discovery"
	"github.com/a754962942/project-common/logs"
	"github.com/a754962942/project-grpc/project"
	"github.com/a754962942/project-project/config"
	project_service_v1 "github.com/a754962942/project-project/pkg/service/project.service.v1"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
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
			project.RegisterRegisterServiceServer(g, project_service_v1.New())
		},
	}
	s := grpc.NewServer()
	g.RegisterFunc(s)
	listen, err := net.Listen("tcp", g.Addr)
	if err != nil {
		logs.LG.Info("connot listen")
	}
	go func() {
		log.Printf("grpc server startd as:%s\n", g.Addr)
		err := s.Serve(listen)
		if err != nil {
			logs.LG.Error(fmt.Sprintf("Server started error %s", err))
			return
		}
	}()
	return s
}
func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.C.Etcd.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	info := discovery.Server{
		Name:    config.C.Gc.Name,
		Addr:    config.C.Gc.Addr,
		Version: config.C.Gc.Version,
		Weight:  config.C.Gc.Weight,
	}
	r := discovery.NewRegister(config.C.Etcd.Addrs, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}
