// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	comment "go-zeroTiktok/system-api/internal/handler/comment"
	favorite "go-zeroTiktok/system-api/internal/handler/favorite"
	feed "go-zeroTiktok/system-api/internal/handler/feed"
	publish "go-zeroTiktok/system-api/internal/handler/publish"
	relation "go-zeroTiktok/system-api/internal/handler/relation"
	user "go-zeroTiktok/system-api/internal/handler/user"
	"go-zeroTiktok/system-api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/tiktok/user/register",
				Handler: user.RegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/tiktok/user/login",
				Handler: user.LoginHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/tiktok/user",
				Handler: user.GetUserMsgHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/tiktok/feed",
				Handler: feed.GetVideoListHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/tiktok/publish/action",
				Handler: publish.ActionHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/tiktok/publish/list",
				Handler: publish.ListHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/tiktok/comment/action",
				Handler: comment.CommentActionHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/tiktok/comment/list/:vedio_id",
				Handler: comment.CommentListHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/tiktok/favorite/action",
				Handler: favorite.FavoriteActionHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/tiktok/favorite/list/:user_id",
				Handler: favorite.FavoriteListHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/tiktok/relation/action",
				Handler: relation.RelationActionHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/tiktok/relation/follow/list/:user_id",
				Handler: relation.RelationFollowListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/tiktok/relation/follower/list/:user_id",
				Handler: relation.RelationFollowerListHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
