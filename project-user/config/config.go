package config

import (
	"github.com/a754962942/project-common/logs"
	"github.com/spf13/viper"
	"log"
	"os"
)

var C = InitConfig()

type Config struct {
	viper       *viper.Viper
	Sc          *ServerConfig
	L           *LogConfig
	R           *RedisConfig
	Gc          *GrpcConfig
	Etcd        *EtcdConfig
	MysqlConfig *MysqlConfig
}
type ServerConfig struct {
	Name string
	Addr string
}
type RedisConfig struct {
	Addr string

	Password string
	DB       int
}
type LogConfig struct {
	DdebugFileName string
	InfoFileName   string
	WarnFileName   string
	MaxSize        int
	MaxAge         int
	MaxBackups     int
}
type EtcdConfig struct {
	Addrs []string
}
type GrpcConfig struct {
	Addr    string
	Name    string
	Version string
	Weight  int64
}
type MysqlConfig struct {
	UserName string
	Password string
	Host     string
	Port     int
	Dbname   string
}

func InitConfig() *Config {
	conf := &Config{viper: viper.New()}
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	pwd, _ := os.Getwd()
	conf.viper.AddConfigPath("/etc/project/user")
	conf.viper.AddConfigPath(pwd + "/config")
	err := conf.viper.ReadInConfig()
	if err != nil {
		logs.LG.Error("viper init failed.")
	}
	conf.readServerConfig()
	conf.readLogConfig()
	conf.readRedisConfig()
	conf.readGrpcConfig()
	conf.readEtcdConfig()
	conf.readMysqlConfig()
	return conf
}

func (c *Config) readServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.Sc = sc
}
func (c *Config) readLogConfig() {
	l := &LogConfig{}
	l.DdebugFileName = c.viper.GetString("zap.debugFileName")
	l.InfoFileName = c.viper.GetString("zap.infoFileName")
	l.WarnFileName = c.viper.GetString("zap.warnFileName")
	l.MaxSize = c.viper.GetInt("maxSize")
	l.MaxAge = c.viper.GetInt("maxAge")
	l.MaxBackups = c.viper.GetInt("maxBackups")
	c.L = l

}
func (c *Config) readRedisConfig() {
	r := &RedisConfig{}
	r.Addr = c.viper.GetString("redis.addr")
	r.Password = c.viper.GetString("redis.password")
	r.DB = c.viper.GetInt("redis.db")
	c.R = r
}
func (c *Config) readGrpcConfig() {
	g := &GrpcConfig{}
	g.Addr = c.viper.GetString("grpc.addr")
	g.Name = c.viper.GetString("grpc.name")
	g.Version = c.viper.GetString("grpc.version")
	g.Weight = c.viper.GetInt64("grpc.weight")
	c.Gc = g
}
func (c *Config) readMysqlConfig() {
	g := &MysqlConfig{}
	g.UserName = c.viper.GetString("mysql.username")
	g.Password = c.viper.GetString("mysql.password")
	g.Host = c.viper.GetString("mysql.host")
	g.Port = c.viper.GetInt("mysql.port")
	g.Dbname = c.viper.GetString("mysql.dbname")
	c.MysqlConfig = g
}
func (c *Config) readEtcdConfig() {
	sc := &EtcdConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	sc.Addrs = addrs
	c.Etcd = sc
}
