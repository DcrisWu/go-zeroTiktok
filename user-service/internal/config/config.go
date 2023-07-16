package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	RedisCfg   redis.RedisConf
	Argon2ID   *Argon2Params
}
