package jwt

import (
	"github.com/ZZINNO/go-zhiliao-common/database/redis"
	"github.com/ZZINNO/go-zhiliao-common/logs"
	"github.com/dgrijalva/jwt-go"
	"sync"
)

var authInstance Auther
var once sync.Once
var jwtInitOpt JwtInitOption

type JwtBaseConfig struct {
	Expire int    //超时,秒
	Key    string //JWT的key
}

type JwtInitOption struct {
	Config    JwtBaseConfig
	RedisSess redis.Session
}

func SetJwtInstanceOption(cfg JwtBaseConfig, sess redis.Session) {
	jwtInitOpt.Config = cfg
	jwtInitOpt.RedisSess = sess
}

func initJWTAuth() {
	var opts []Option
	opts = append(opts, SetExpired(jwtInitOpt.Config.Expire))
	opts = append(opts, SetSigningKey([]byte(jwtInitOpt.Config.Key)))
	opts = append(opts, SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(jwtInitOpt.Config.Key), nil
	}))
	opts = append(opts, SetSigningMethod(jwt.SigningMethodHS512))

	var store Storer
	store, err := NewStore(jwtInitOpt.RedisSess, jwtInitOpt.Config.Key)
	if err != nil {
		logs.Error("初始化JWT错误", err)
	}
	authInstance = New(store, opts...)
}

func GetAuthInstance() Auther {
	once.Do(initJWTAuth)
	return authInstance
}
