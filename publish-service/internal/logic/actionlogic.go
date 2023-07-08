package logic

import (
	"context"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/publish-service/internal/svc"
	"go-zeroTiktok/publish-service/pb/publish"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionLogic {
	return &ActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Action 发布视频接口
func (l *ActionLogic) Action(in *publish.ActionReq) (*publish.ActionResp, error) {
	videoModel := &db.Video{
		AuthorID: int(in.AuthorId),
		PlayUrl:  in.PlayUrl,
		CoverUrl: in.CoverUrl,
		Title:    in.Title,
	}
	err := db.CreateVideo(l.ctx, l.svcCtx.DB, videoModel)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	return &publish.ActionResp{
		IsSuccess: true,
	}, nil
}
