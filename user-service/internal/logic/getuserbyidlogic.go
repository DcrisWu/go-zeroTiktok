package logic

import (
	"context"

	"go-zeroTiktok/user-service/internal/svc"
	"go-zeroTiktok/user-service/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByIdLogic) GetUserById(in *user.UserReq) (*user.UserResp, error) {
	// todo: add your logic here and delete this line

	return &user.UserResp{}, nil
}
