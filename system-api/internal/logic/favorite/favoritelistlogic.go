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

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList(req *types.FavoriteListReq) (*types.FavoriteListResp, error) {
	uid := utils.GetUid(l.ctx)
	isExit, err := utils.IsJwtInRedis(l.ctx, l.svcCtx.Redis, uid)
	if err != nil || !isExit {
		return &types.FavoriteListResp{
			StatusCode: http.StatusUnauthorized,
			StatusMsg:  "请先登录",
		}, nil
	}
	if uid == utils.UidNotFound || req.UserId == "0" {
		return &types.FavoriteListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return &types.FavoriteListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}

	resp, err := l.svcCtx.FavoriteService.GetFavoriteList(l.ctx, &favorite.FavoriteListReq{
		UserId: userId,
	})
	if err != nil {
		return &types.FavoriteListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	var videos []*types.Video
	for _, v := range resp.VideoList {
		videos = append(videos, &types.Video{
			Id: v.Id,
			Author: &types.User{
				Id:            v.Author.Id,
				Name:          v.Author.Name,
				FollowCount:   v.Author.FollowCount,
				FollowerCount: v.Author.FollowerCount,
				IsFollow:      v.Author.IsFollow,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		})
	}
	return &types.FavoriteListResp{
		StatusCode: utils.SUCCESS,
		StatusMsg:  "成功",
		VideoList:  videos,
	}, nil
}
