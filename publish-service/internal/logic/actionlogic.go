package logic

import (
	"context"
	"go-zeroTiktok/model"
	"go-zeroTiktok/publish-service/internal/svc"
	"go-zeroTiktok/publish-service/pb/publish"
	"go-zeroTiktok/utils"
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

func (l *ActionLogic) Action(in *publish.ActionReq) (*publish.ActionResp, error) {
	vid, err := utils.NewBasicGenerator().GenerateId()
	if err != nil {
		return nil, status.Error(500, "id生成失败")
	}
	_, err = l.svcCtx.VideoModel.Insert(l.ctx, &model.Video{
		Id:       vid,
		AuthorId: in.AuthorId,
		PlayUrl:  in.PlayUrl,
		CoverUrl: in.CoverUrl,
		Title:    in.Title,
	})
	if err != nil {
		return nil, status.Error(500, "发布失败")
	}

	return &publish.ActionResp{
		VideoId: vid,
	}, nil
}
