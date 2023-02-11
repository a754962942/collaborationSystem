package user

import (
	loginServiceV1 "github.com/a754962942/project-user/pkg/service/login.service.v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var LoginServiceClient loginServiceV1.LoginServiceClient

func InitRpcUserClient() {
	conn, err := grpc.Dial("127.0.0.1:8881", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect:%v\n", err)
	}
	LoginServiceClient = loginServiceV1.NewLoginServiceClient(conn)
}
