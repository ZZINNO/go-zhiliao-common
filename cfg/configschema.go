package cfg

import (
	"github.com/ZZINNO/go-zhiliao-common/database/mysql"
	"github.com/ZZINNO/go-zhiliao-common/database/redis"
	"github.com/ZZINNO/go-zhiliao-common/lib/jwt"
	"github.com/spf13/viper"
)

const (
	DEFAULT_CONFIG_PATH          = "./setting"
	DEFAULT_CONFIG_PATH_MULTIAPP = "../../setting"
	DEFAULT_CONFIG_FILE          = "debug.config"
)

//只提供最基础的配置结构体
type BaseConfig struct {
	AppConfig
}

type AppConfig struct {
	Project   string
	Domain    string
	SentryUrl string
	Env       string
	Jwt       jwt.JwtBaseConfig
	Mysql     map[string]mysql.MysqlBaseConfig
	Redis     redis.RedisBaseConfig
}

type ConfParser struct {
	Path             []string
	ConfigName       string
	OnConfigChangeCb func()
	viper            *viper.Viper
}
