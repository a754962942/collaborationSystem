package gorms

import (
	"fmt"
	"github.com/a754962942/project-user/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _db *gorm.DB

func init() {
	//配置MySQL连接参数
	username := config.C.MysqlConfig.UserName //账号
	password := config.C.MysqlConfig.Password //密码
	host := config.C.MysqlConfig.Host         //数据库地址，可以是Ip或者域名
	port := config.C.MysqlConfig.Port         //数据库端口
	Dbname := config.C.MysqlConfig.Dbname     //数据库名
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	var err error
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
}
func GetDB() *gorm.DB {
	return _db
}
