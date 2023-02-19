package project

import (
	"github.com/a754962942/project-api/config"
	"github.com/a754962942/project-common/discovery"
	"github.com/a754962942/project-common/logs"
	"github.com/a754962942/project-grpc/project"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
)

var RegisterServiceClient project.RegisterServiceClient

func InitRpcUserClient() {
	etcdRegister := discovery.NewResolver(config.C.Etcd.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd:///project", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect:%v\n", err)
	}
	RegisterServiceClient = project.NewRegisterServiceClient(conn)
}
