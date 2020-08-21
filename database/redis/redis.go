package redis

import (
	"github.com/go-redis/redis"
	"strings"
	"time"
)

type Session struct {
	Init             bool //是否已经完成了初始化
	mc               *redis.Client
	expire           time.Duration
	defaultKeyPrefix string
}

type RedisBaseConfig struct {
	Sentinel  bool
	Host      string
	Pass      string
	Master    string
	DB        int
	KeyPrefix string
}

type RedisCustomConfig struct {
	DB        int
	KeyPrefix string
}

var DEFAULT_EXPIRE = 10 * time.Minute

func CombineConfig(r RedisBaseConfig, i RedisCustomConfig) RedisBaseConfig {
	o := r
	o.DB = i.DB
	o.KeyPrefix = i.KeyPrefix
	return o
}

func (sess *Session) GetMc() *redis.Client {
	return sess.mc
}

//初始化
//默认过期时间是10分钟
func InitSession(cfg RedisBaseConfig) (*Session, error) {
	sess := Session{}
	c, err := InitRedisClient(cfg)
	if err != nil {
		return &sess, err
	}
	sess.mc = c
	sess.expire = DEFAULT_EXPIRE
	sess.defaultKeyPrefix = cfg.KeyPrefix
	sess.Init = true
	return &sess, nil
}

//默认key前缀
//最后拼装的key名为：项目名_模块名_{name}  name是你传入的名称
func (sess *Session) GetWithDefaultKey(keyName string) string {
	return sess.defaultKeyPrefix + keyName
}

func (sess *Session) Close() error {
	return sess.mc.Close()
}

//自增某一个key
func (sess *Session) Incr(key string) {
	if len(key) > 0 {
		sess.mc.Incr(key)
	}
}

//为某个key设置过期时间
func (sess *Session) Expire(key string, dur time.Duration) {
	if len(key) > 0 {
		sess.mc.Expire(key, dur)
	}
}

//取出k-v
func (sess *Session) Get(key string) string {
	return sess.mc.Get(key).Val()
}

//按照pattern获取匹配的key列表
func (sess *Session) GetPerfix(key string) []string {
	return sess.mc.Keys(key).Val()
}

//存放k-v并且设置10分钟过期
func (sess *Session) Set(key string, val string) {
	if len(key) > 0 {
		sess.mc.Set(key, val, sess.expire)
	}
}

func (sess *Session) SetWithExpire(key, val string, expiration time.Duration) error {
	cmd := sess.mc.Set(key, val, expiration)
	return cmd.Err()
}

//删除key
func (sess *Session) Delete(key string) error {
	if len(key) > 0 {
		return sess.mc.Del(key).Err()
	}
	return nil
}

func (sess *Session) Exists(keys string) (bool, error) {
	cmd := sess.mc.Exists(keys)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

//删除并返回结果
func (sess *Session) DeleteAndReturn(key string) string {
	if len(key) > 0 {
		pipe := sess.mc.Pipeline()
		pipe.Get(key).Val()
		pipe.Del(key)
		cmder, _ := pipe.Exec()
		if len(cmder) > 0 {
			strmap := cmder[0].(*redis.StringCmd).Val()
			return strmap
		}
	}
	return ""
}

//队列入队
func (sess *Session) Rpush(queue string, obj string) {
	sess.mc.RPush(queue, obj)
}

//队列出队并返回
func (sess *Session) Rpop(queue string) string {
	return sess.mc.RPop(queue).Val()
}

//判断k-v是否相等
func (sess *Session) Check(key string, val string) bool {
	if sess.mc.Get(key).Val() == val {
		return true
	} else {
		return false
	}
}

//获取队列长度
func (sess *Session) LLen(key string) int64 {
	return sess.mc.LLen(key).Val()
}

func InitRedisClient(cfg RedisBaseConfig) (client *redis.Client, err error) {
	if cfg.Sentinel {
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: strings.Split(cfg.Host, ","),
			Password:      cfg.Pass,
			MasterName:    cfg.Master,
			DB:            cfg.DB,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Host,
			Password: cfg.Pass,
			DB:       cfg.DB,
		})
	}
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
