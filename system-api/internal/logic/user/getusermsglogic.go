package user

import (
	"context"

	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMsgLogic {
	return &GetUserMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserMsgLogic) GetUserMsg(req *types.UserReq) (resp *types.UserResp, err error) {
	// todo: add your logic here and delete this line

	return
}
