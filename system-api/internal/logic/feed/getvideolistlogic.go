package feed

import (
	"context"

	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *GetVideoListLogic) GetVideoList(req *types.FeedReq) (resp *types.FeedResp, err error) {
	// todo: add your logic here and delete this line

	return
}
