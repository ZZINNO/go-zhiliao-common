package jwt

import (
	"errors"
	"fmt"
	"github.com/ZZINNO/go-zhiliao-common/database/redis"
	"time"
)

var ERR_SESS_NotInit = "redis session还没初始化"

// Storer 黑名单存储接口
type Storer interface {
	// 放入令牌，指定到期时间
	Set(tokenString string, expiration time.Duration) error
	// 检查令牌是否存在
	Exists(tokenString string) (bool, error)
	// 关闭存储
	Close() error
}

func NewStore(sess redis.Session, prefix string) (*Store, error) {
	if sess.Init {
		return &Store{
			cli:    &sess,
			prefix: prefix,
		}, nil
	} else {
		return &Store{}, errors.New(ERR_SESS_NotInit)
	}
}

type redisClient interface {
	Get(key string) string
	SetWithExpire(key, val string, expiration time.Duration) error
	Exists(key string) (bool, error)
	Close() error
}

// Store redis存储
type Store struct {
	cli    redisClient
	prefix string
}

func (a *Store) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", a.prefix, key)
}

// Set ...
func (a *Store) Set(tokenString string, expiration time.Duration) error {
	return a.cli.SetWithExpire(a.wrapperKey(tokenString), "1", expiration)
}

// Exists ...
func (a *Store) Exists(tokenString string) (bool, error) {
	return a.cli.Exists(a.wrapperKey(tokenString))
}

// Close ...
func (a *Store) Close() error {
	return a.cli.Close()
}
