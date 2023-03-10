package main

import (
	_ "github.com/a754962942/project-api/api"
	"github.com/a754962942/project-api/config"
	"github.com/a754962942/project-api/router"
	common "github.com/a754962942/project-common"
	"github.com/a754962942/project-common/logs"
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
	common.Run(r, config.C.Sc.Name, config.C.Sc.Addr, nil)
}
