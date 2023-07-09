package logic

import (
	"context"

	"go-zeroTiktok/favorite-service/internal/svc"
	"go-zeroTiktok/favorite-service/pb/favorite"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelFavoriteLogic {
	return &CancelFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelFavoriteLogic) CancelFavorite(in *favorite.CancelFavoriteReq) (*favorite.CancelFavoriteResp, error) {
	// todo: add your logic here and delete this line

	return &favorite.CancelFavoriteResp{}, nil
}
