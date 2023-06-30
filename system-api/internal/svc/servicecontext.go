package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zeroTiktok/system-api/internal/config"
	"go-zeroTiktok/user-service/pb/user"
)

type ServiceContext struct {
	Config      config.Config
	UserService user.UserClient
}

var cfg = zrpc.RpcClientConf{
	Endpoints: []string{"127.0.0.1:8080"},
}

func NewServiceContext(c config.Config) *ServiceContext {

	conn := zrpc.MustNewClient(cfg)
	return &ServiceContext{
		Config:      c,
		UserService: user.NewUserClient(conn.Conn()),
	}
}
