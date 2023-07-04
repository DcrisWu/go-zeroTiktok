package publish

import (
	"context"
	"go-zeroTiktok/publish-service/pb/publish"
	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"
	"go-zeroTiktok/utils"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.PublishListReq) (*types.PublishListResp, error) {
	var userId int64
	// 假如没有传uid,就默认是自己
	if req.UserId == 0 {
		userId = utils.GetUid(l.ctx)
	} else {
		userId = req.UserId
	}
	list, err := l.svcCtx.PublishService.List(l.ctx, &publish.ListReq{
		AuthorId: userId,
	})
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	videoList := make([]*types.Video, 0)
	for _, i := range list.VideoList {
		videoList = append(videoList, &types.Video{
			Id: i.Id,
			Author: &types.User{
				Id:            i.Author.Id,
				Name:          i.Author.Name,
				FollowCount:   i.Author.FollowCount,
				FollowerCount: i.Author.FollowerCount,
				IsFollow:      i.Author.IsFollow,
			},
			PlayUrl:       i.PlayUrl,
			CoverUrl:      i.CoverUrl,
			FavoriteCount: i.FavorCount,
			CommentCount:  i.CommentCount,
			IsFavorite:    i.IsFavor,
			Title:         i.Title,
		})
	}

	return &types.PublishListResp{
		StatusCode: utils.SUCCESS,
		VideoList:  videoList,
	}, nil
}
