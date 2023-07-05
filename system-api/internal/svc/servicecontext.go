package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zeroTiktok/feed-service/pb/feed"
	"go-zeroTiktok/publish-service/pb/publish"
	"go-zeroTiktok/system-api/internal/config"
	"go-zeroTiktok/user-service/pb/user"
)

type ServiceContext struct {
	Config         config.Config
	UserService    user.UserClient
	PublishService publish.PublishClient
	FeedService    feed.FeedClient
}

var (
	userCfg = zrpc.RpcClientConf{
		Endpoints: []string{"127.0.0.1:8080"},
	}
	publishCfg = zrpc.RpcClientConf{
		Endpoints: []string{"127.0.0.1:8081"},
	}
	feedCfg = zrpc.RpcClientConf{
		Endpoints: []string{"127.0.0.1:8082"},
	}
)

func NewServiceContext(c config.Config) *ServiceContext {

	userConn := zrpc.MustNewClient(userCfg)
	publishConn := zrpc.MustNewClient(publishCfg)
	feedConn := zrpc.MustNewClient(feedCfg)
	return &ServiceContext{
		Config:         c,
		UserService:    user.NewUserClient(userConn.Conn()),
		PublishService: publish.NewPublishClient(publishConn.Conn()),
		FeedService:    feed.NewFeedClient(feedConn.Conn()),
	}
}
