package svc

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zeroTiktok/comment-service/pb/comment"
	"go-zeroTiktok/favorite-service/pb/favorite"
	"go-zeroTiktok/feed-service/pb/feed"
	"go-zeroTiktok/publish-service/pb/publish"
	"go-zeroTiktok/relation-service/pb/relation"
	"go-zeroTiktok/system-api/internal/config"
	"go-zeroTiktok/user-service/pb/user"
)

type ServiceContext struct {
	Config          config.Config
	Redis           *redis.Redis
	UserService     user.UserClient
	PublishService  publish.PublishClient
	FeedService     feed.FeedClient
	CommentService  comment.CommentClient
	FavoriteService favorite.FavoriteClient
	RelationService relation.RelationClient
}

//var (
//	userCfg = zrpc.RpcClientConf{
//		Endpoints: []string{"127.0.0.1:8080"},
//	}
//	publishCfg = zrpc.RpcClientConf{
//		Endpoints: []string{"127.0.0.1:8081"},
//	}
//	feedCfg = zrpc.RpcClientConf{
//		Endpoints: []string{"127.0.0.1:8082"},
//	}
//	commentCfg = zrpc.RpcClientConf{
//		Endpoints: []string{"127.0.0.1:8083"},
//	}
//	favoriteCfg = zrpc.RpcClientConf{
//		Endpoints: []string{"127.0.0.1:8084"},
//	}
//	relationCfg = zrpc.RpcClientConf{
//		Endpoints: []string{"127.0.0.1:8085"},
//	}
//)

func NewServiceContext(c config.Config) *ServiceContext {
	userConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: c.Etcd.Hosts,
			Key:   c.Etcd.UserKey,
		},
	})
	publishConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: c.Etcd.Hosts,
			Key:   c.Etcd.PublishKey,
		},
	})
	feedConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: c.Etcd.Hosts,
			Key:   c.Etcd.FeedKey,
		},
	})
	commentConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: c.Etcd.Hosts,
			Key:   c.Etcd.CommentKey,
		},
	})
	favoriteConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: c.Etcd.Hosts,
			Key:   c.Etcd.FavoriteKey,
		},
	})
	relationConn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: c.Etcd.Hosts,
			Key:   c.Etcd.RelationKey,
		},
	})
	//userConn := zrpc.MustNewClient(userCfg)
	//publishConn := zrpc.MustNewClient(publishCfg)
	//feedConn := zrpc.MustNewClient(feedCfg)
	//commentConn := zrpc.MustNewClient(commentCfg)
	//favoriteConn := zrpc.MustNewClient(favoriteCfg)
	//relationConn := zrpc.MustNewClient(relationCfg)
	return &ServiceContext{
		Config:          c,
		Redis:           redis.MustNewRedis(c.RedisCfg),
		UserService:     user.NewUserClient(userConn.Conn()),
		PublishService:  publish.NewPublishClient(publishConn.Conn()),
		FeedService:     feed.NewFeedClient(feedConn.Conn()),
		CommentService:  comment.NewCommentClient(commentConn.Conn()),
		FavoriteService: favorite.NewFavoriteClient(favoriteConn.Conn()),
		RelationService: relation.NewRelationClient(relationConn.Conn()),
	}
}
