package logic

import (
	"context"
	"go-zeroTiktok/model"
	"google.golang.org/grpc/status"

	"go-zeroTiktok/publish-service/internal/svc"
	"go-zeroTiktok/publish-service/pb/publish"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLogic) List(in *publish.ListReq) (*publish.ListResp, error) {
	list, err := l.svcCtx.VideoModel.GetVideosByAuthorId(l.ctx, in.AuthorId, 10, 0)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	var videoList []*publish.Video
	for _, i := range list {
		one, err := l.svcCtx.UserModel.FindOne(l.ctx, i.AuthorId)
		if err != nil {
			return nil, status.Error(500, err.Error())
		}
		isFollow := false
		if in.Uid != in.AuthorId {
			_, err := l.svcCtx.RelationModel.FindOneByUserIdToUserId(l.ctx, in.Uid, i.AuthorId)
			if err != nil && err != model.ErrNotFound {
				return nil, status.Error(500, err.Error())
			}
			if err == nil {
				isFollow = true
			}
		}
		_, err = l.svcCtx.FavoriteModel.FindOneByUserIdVideoId(l.ctx, in.Uid, i.Id)
		if err != nil && err != model.ErrNotFound {
			return nil, status.Error(500, err.Error())
		}
		isFavorite := false
		if err == nil {
			isFavorite = true
		}
		videoList = append(videoList, &publish.Video{
			Id: i.Id,
			Author: &publish.User{
				Id:            one.Id,
				Name:          one.UserName,
				FollowCount:   one.FollowingCount,
				FollowerCount: one.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:      i.PlayUrl,
			CoverUrl:     i.CoverUrl,
			FavorCount:   i.FavoriteCount,
			CommentCount: i.CommentCount,
			IsFavor:      isFavorite,
			Title:        i.Title,
		})
	}
	return &publish.ListResp{
		VideoList: videoList,
	}, nil
}
