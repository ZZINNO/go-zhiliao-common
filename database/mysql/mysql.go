package mysql

import (
	"github.com/ZZINNO/go-zhiliao-common/logs"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
	"github.com/xormplus/xorm/log"
)

type MysqlBaseConfig struct {
	DSN        string //如果dsn存在的情况下,则下面不生效,否则将会自动拼接
	ShowSql    bool
	LogLevel   string
	Net        string
	User       string
	Pass       string
	Addr       string
	DBName     string
	NativePass bool
	Params     map[string]string
}

func GetDataSource(cfg MysqlBaseConfig) string {
	c := mysql.Config{
		Net:                  cfg.Net,
		User:                 cfg.User,
		Passwd:               cfg.Pass,
		Addr:                 cfg.Addr,
		DBName:               cfg.DBName,
		Params:               cfg.Params,
		AllowNativePasswords: cfg.NativePass,
	}
	if cfg.DSN != "" {
		return cfg.DSN
	} else {
		return c.FormatDSN()
	}
}

func InitMysqlEngine(cfg MysqlBaseConfig) (engine *xorm.Engine, err error) {
	dsn := GetDataSource(cfg)
	DBEngine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		logs.Error("Mysql连接失败", err)
		return nil, err
	}
	if err := DBEngine.Ping(); err != nil {
		return nil, err
	}
	DBEngine.SetLogger(log.NewSimpleLogger(logs.Logger.GetWriter()))
	DBEngine.ShowSQL(cfg.ShowSql)
	l := log.LOG_DEBUG
	switch cfg.LogLevel {
	case "debug":
		l = log.LOG_DEBUG
	case "info":
		l = log.LOG_INFO
	case "warning":
		l = log.LOG_WARNING
	default:
		l = log.LOG_DEBUG
	}
	DBEngine.Logger().SetLevel(l)
	//这里需要设定mapper不然会在sql的时候把id拼接成I_D
	DBEngine.SetMapper(core.SameMapper{})
	//链接池
	DBEngine.SetMaxIdleConns(25)
	DBEngine.SetMaxOpenConns(256)
	return DBEngine, nil
}
