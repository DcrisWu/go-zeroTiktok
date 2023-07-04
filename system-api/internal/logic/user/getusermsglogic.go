package user

import (
	"context"
	"go-zeroTiktok/user-service/pb/user"
	"go-zeroTiktok/utils"
	"google.golang.org/grpc/status"

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
	uid := utils.GetUid(l.ctx)
	if uid == utils.UidNotFound {
		return &types.UserResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "用户不存在",
		}, nil
	}
	userResp, err := l.svcCtx.UserService.GetUserById(l.ctx, &user.UserReq{
		Uid:    uid,
		UserId: req.UserId,
	})
	if err != nil {
		if s, ok := status.FromError(err); !ok {
			return &types.UserResp{
				StatusCode: utils.FAILED,
				StatusMsg:  s.Message(),
			}, nil
		} else {
			return nil, err
		}
	}
	return &types.UserResp{
		StatusCode: utils.SUCCESS,
		StatusMsg:  "获取用户信息成功",
		User: &types.User{
			Id:            userResp.User.Id,
			Name:          userResp.User.Name,
			FollowCount:   userResp.User.FollowCount,
			FollowerCount: userResp.User.FollowerCount,
			IsFollow:      userResp.IsFollow,
		},
	}, nil
}
