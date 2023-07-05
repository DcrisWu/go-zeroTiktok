package feed

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/feed-service/pb/feed"
	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"
	"go-zeroTiktok/utils"
	"strconv"
	"strings"
)

type GetVideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListLogic {
	return &GetVideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVideoListLogic) GetVideoList(req *types.FeedReq) (*types.FeedResp, error) {
	var uid = utils.UidNotFound
	if req.Authorization != "" {
		strArr := strings.Split(req.Authorization, " ")
		if len(strArr) != 2 {
			goto flag
		}
		token := strArr[1]
		jwt, err := utils.ParseJWT(token, l.svcCtx.Config.Auth.AccessSecret)
		if err != nil {
			goto flag
		}
		parseInt, err := strconv.ParseInt(jwt["uid"].(string), 10, 64)
		if err != nil {
			goto flag
		}
		uid = parseInt
	}

flag:
	in := &feed.FeedReq{
		Uid:      uid,
		LastTime: req.LatestTime,
	}
	resp, err := l.svcCtx.FeedService.Feed(l.ctx, in)
	if err != nil {
		return &types.FeedResp{
			StatsCode: utils.FAILED,
			StatusMsg: err.Error(),
		}, nil
	}
	videoList := make([]*types.Video, 0)
	for _, v := range resp.VideoList {
		videoList = append(videoList, &types.Video{
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
			FavoriteCount: v.FavorCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavor,
			Title:         v.Title,
		})
	}
	return &types.FeedResp{
		StatsCode: utils.SUCCESS,
		StatusMsg: "获取成功",
		NextTime:  resp.EarliestTime,
		VideoList: videoList,
	}, nil
}
