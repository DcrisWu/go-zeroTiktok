package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/feed-service/internal/svc"
	"go-zeroTiktok/feed-service/pb/feed"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/models/pack"
	"google.golang.org/grpc/status"
	"time"
)

type FeedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

const (
	LIMIT = 30
)

func (l *FeedLogic) Feed(in *feed.FeedReq) (*feed.FeedResp, error) {
	vis, nextTime, err := l.getUserFeed(in, in.Uid)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	return &feed.FeedResp{
		VideoList: vis,
		NextTime:  nextTime,
	}, nil
}

func (l *FeedLogic) getUserFeed(req *feed.FeedReq, fromId int64) (vis []*feed.Video, nextTime int64, err error) {
	videos, err := db.GetVideos(l.ctx, l.svcCtx.DB, LIMIT, req.LastTime)
	if err != nil {
		return nil, nextTime, err
	}
	if len(videos) == 0 {
		nextTime = time.Now().UnixMilli()
		return vis, nextTime, nil
	} else {
		nextTime = videos[len(videos)-1].UpdatedAt.UnixMilli()
	}
	if vis, err = pack.FeedVideos(l.ctx, l.svcCtx.DB, videos, fromId); err != nil {
		nextTime = time.Now().UnixMilli()
		return nil, nextTime, err
	}
	return vis, nextTime, nil
}

//func (l *FeedLogic) Feed(in *feed.FeedReq) (*feed.FeedResp, error) {
//	uid := in.Uid
//	lastTime := mq.TimeStringToGoTime(in.LastTime)
//	if lastTime == (time.Time{}) {
//		lastTime = time.Now()
//	}
//	list, err := l.svcCtx.VideoModel.GetVideosOrderByTime(l.ctx, 30, 0, lastTime)
//	if err != nil {
//		return nil, status.Error(500, err.Error())
//	}
//
//	var earliestTime = time.Now()
//	var videoList []*feed.FeedVideo
//	for _, i := range list {
//		one, err := l.svcCtx.UserModel.FindOne(l.ctx, i.AuthorId)
//		earliestTime = minTime(earliestTime, i.CreatedAt)
//		if err != nil {
//			return nil, status.Error(500, err.Error())
//		}
//		// 如果携带了token访问，就判断是否点赞了视频，是否关注了作者
//		isFollow := false
//		isFavorite := false
//		if uid != mq.UidNotFound {
//			_, err := l.svcCtx.RelationModel.FindOneByUserIdToUserId(l.ctx, uid, i.AuthorId)
//			if err != nil && err != model.ErrNotFound {
//				return nil, status.Error(500, err.Error())
//			}
//			if err == nil {
//				isFollow = true
//			}
//
//			_, err = l.svcCtx.FavoriteModel.FindOneByUserIdVideoId(l.ctx, uid, i.Id)
//			if err != nil && err != model.ErrNotFound {
//				return nil, status.Error(500, err.Error())
//			}
//			if err == nil {
//				isFavorite = true
//			}
//		}
//		videoList = append(videoList, &feed.FeedVideo{
//			Id: i.Id,
//			Author: &feed.User{
//				Id:            one.Id,
//				Name:          one.UserName,
//				FollowCount:   one.FollowingCount,
//				FollowerCount: one.FollowerCount,
//				IsFollow:      isFollow,
//			},
//			PlayUrl:      i.PlayUrl,
//			CoverUrl:     i.CoverUrl,
//			FavorCount:   i.FavoriteCount,
//			CommentCount: i.CommentCount,
//			IsFavor:      isFavorite,
//			Title:        i.Title,
//		})
//	}
//	return &feed.FeedResp{
//		VideoList:    videoList,
//		EarliestTime: earliestTime.Unix(),
//	}, nil
//}
//
//func minTime(a, b time.Time) time.Time {
//	if a.Before(b) {
//		return a
//	}
//	return b
//}
