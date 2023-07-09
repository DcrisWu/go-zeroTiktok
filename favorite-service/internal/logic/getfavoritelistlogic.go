package logic

import (
	"context"

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

func (l *GetFavoriteListLogic) GetFavoriteList(in *favorite.FavoriteList) (*favorite.FavoriteListResp, error) {
	// todo: add your logic here and delete this line

	return &favorite.FavoriteListResp{}, nil
}
