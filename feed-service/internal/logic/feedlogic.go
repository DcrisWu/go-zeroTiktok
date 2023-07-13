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
