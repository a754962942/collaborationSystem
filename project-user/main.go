package main

import (
	common "github.com/a754962942/project-common"
	"github.com/a754962942/project-common/logs"
	"github.com/a754962942/project-user/config"
	"github.com/a754962942/project-user/router"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	lc := &logs.LogConfig{
		DebugFileName: config.C.L.DdebugFileName,
		InfoFileName:  config.C.L.InfoFileName,
		WarnFileName:  config.C.L.WarnFileName,
		MaxSize:       config.C.L.MaxSize,
		MaxAge:        config.C.L.MaxAge,
		MaxBackups:    config.C.L.MaxBackups,
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
	router.InitRouter(r)
	grpc := router.RegisterGrpc()
	stop := func() {
		grpc.Stop()
	}
	common.Run(r, config.C.Sc.Name, config.C.Sc.Addr, stop)
}
