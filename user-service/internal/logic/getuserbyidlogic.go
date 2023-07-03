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
	one, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	if one == nil {
		return &user.UserResp{
			Status: 1,
		}, nil
	}

	return &user.UserResp{
		Status: 0,
		User: &user.UserInfo{
			Id:            one.Id,
			Name:          one.UserName,
			FollowCount:   one.FollowingCount,
			FollowerCount: one.FollowerCount,
		},
		// todo isFollow
	}, nil
}
