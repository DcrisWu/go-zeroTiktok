package logic

import (
	"context"
	"go-zeroTiktok/model"
	"go-zeroTiktok/user-service/internal/svc"
	"go-zeroTiktok/user-service/pb/user"
	"google.golang.org/grpc/status"

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
	if in.Uid == 0 || in.UserId == 0 {
		return nil, status.Error(400, "参数缺失")
	}
	one, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		logx.Error(err)
		return nil, status.Error(500, err.Error())
	}

	if one == nil {
		return &user.UserResp{
			Status: 1,
		}, nil
	}

	isFollow := false
	if in.UserId != in.Uid {
		rel, err := l.svcCtx.RelationModel.FindOneByUserIdToUserId(l.ctx, in.Uid, in.UserId)
		if err != nil && err != model.ErrNotFound {
			return nil, status.Error(500, err.Error())
		}
		if rel != nil {
			isFollow = true
		}
	}

	return &user.UserResp{
		Status: 0,
		User: &user.UserInfo{
			Id:            one.Id,
			Name:          one.UserName,
			FollowCount:   one.FollowingCount,
			FollowerCount: one.FollowerCount,
		},
		IsFollow: isFollow,
	}, nil
}
