package logic

import (
	"context"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/models/pack"
	"google.golang.org/grpc/status"

	"go-zeroTiktok/favorite-service/internal/svc"
	"go-zeroTiktok/favorite-service/pb/favorite"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteListLogic {
	return &GetFavoriteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteListLogic) GetFavoriteList(in *favorite.FavoriteListReq) (*favorite.FavoriteListResp, error) {
	list, err := db.FavoriteList(l.ctx, l.svcCtx.DB, in.UserId)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	videos, err := pack.FavoriteVideos(l.ctx, l.svcCtx.DB, list, in.UserId)
	if err != nil {
		return nil, err
	}
	return &favorite.FavoriteListResp{
		VideoList: videos,
	}, nil
}
