package favorite

import (
	"context"
	"go-zeroTiktok/favorite-service/pb/favorite"
	"go-zeroTiktok/utils"
	"net/http"
	"strconv"

	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteActionLogic) FavoriteAction(req *types.FavoriteActionReq) (*types.FavoriteActionResp, error) {
	uid := utils.GetUid(l.ctx)
	isExit, err := utils.IsJwtInRedis(l.ctx, l.svcCtx.Redis, uid)
	if err != nil || !isExit {
		return &types.FavoriteActionResp{
			StatusCode: http.StatusUnauthorized,
			StatusMsg:  "请先登录",
		}, nil
	}
	if (req.ActionType != "1" && req.ActionType != "2") || uid == utils.UidNotFound || req.VideoId == "0" {
		return &types.FavoriteActionResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	vid, err := strconv.ParseInt(req.VideoId, 10, 64)
	if err != nil {
		return &types.FavoriteActionResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	_, err = l.svcCtx.FavoriteService.FavoriteAction(l.ctx, &favorite.FavoriteActionReq{
		UserId:  uid,
		VideoId: vid,
	})
	if err != nil {
		return &types.FavoriteActionResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	return &types.FavoriteActionResp{
		StatusCode: utils.SUCCESS,
		StatusMsg:  "成功",
	}, nil
}
