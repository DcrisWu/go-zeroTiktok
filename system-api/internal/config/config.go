package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type EtcdCfg struct {
	Hosts       []string
	UserKey     string
	PublishKey  string
	FeedKey     string
	CommentKey  string
	FavoriteKey string
	RelationKey string
}

type Config struct {
	rest.RestConf
	RedisCfg redis.RedisConf
	Auth     struct {
		AccessSecret string
		AccessExpire int64
	}
	Etcd EtcdCfg
}
