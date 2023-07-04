package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zeroTiktok/publish-service/pb/publish"
	"go-zeroTiktok/system-api/internal/config"
	"go-zeroTiktok/user-service/pb/user"
)

type ServiceContext struct {
	Config         config.Config
	UserService    user.UserClient
	PublishService publish.PublishClient
}

var (
	userCfg = zrpc.RpcClientConf{
		Endpoints: []string{"127.0.0.1:8080"},
	}
	publishCfg = zrpc.RpcClientConf{
		Endpoints: []string{"127.0.0.1:8081"},
	}
)

func NewServiceContext(c config.Config) *ServiceContext {

	userConn := zrpc.MustNewClient(userCfg)
	publishConn := zrpc.MustNewClient(publishCfg)
	return &ServiceContext{
		Config:         c,
		UserService:    user.NewUserClient(userConn.Conn()),
		PublishService: publish.NewPublishClient(publishConn.Conn()),
	}
}
