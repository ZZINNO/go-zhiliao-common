package network

import (
	"github.com/ZZINNO/go-zhiliao-common/lib/daemon"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpBaseConfig struct {
	Addr     string `json:"addr"`
	GlaceFul int
}

type httpApiServerOption struct {
	Address string
	Env     string
	R       *gin.Engine
}

type HttpApiServerOption struct {
	f func(option *httpApiServerOption)
}

func SetHttpAddress(addr string) HttpApiServerOption {
	return HttpApiServerOption{func(option *httpApiServerOption) {
		option.Address = addr
	}}
}

func SetHttpEngine(r *gin.Engine) HttpApiServerOption {
	return HttpApiServerOption{func(option *httpApiServerOption) {
		option.R = r
	}}
}

func SetHttpEnv(env string) HttpApiServerOption {
	return HttpApiServerOption{func(option *httpApiServerOption) {
		option.Env = env
	}}
}
func InitApiServer(opts ...HttpApiServerOption) error {
	//设置默认值
	opt := httpApiServerOption{
		Address: "0.0.0.0:8000",
		R:       &gin.Engine{},
		Env:     "debug",
	}
	for _, optFunc := range opts {
		optFunc.f(&opt)
	}
	if opt.Env == "product" {
		gin.SetMode("release")
	} else {
		gin.SetMode("debug")
		opt.R.Use(gin.Logger())
	}
	return opt.R.Run(opt.Address)
}

func InitApiServerDaemon(opts ...HttpApiServerOption) error {
	//设置默认值
	opt := httpApiServerOption{
		Address: "0.0.0.0:8000",
		R:       &gin.Engine{},
		Env:     "debug",
	}
	for _, optFunc := range opts {
		optFunc.f(&opt)
	}
	if opt.Env == "product" {
		gin.SetMode("release")
	} else {
		gin.SetMode("debug")
		opt.R.Use(gin.Logger())
	}
	httpSrv := http.Server{
		Addr:    opt.Address,
		Handler: opt.R,
	}
	daemon.DaemonCtl.ZeroDown(httpSrv)
	return nil
}
