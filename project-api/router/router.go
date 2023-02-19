package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
