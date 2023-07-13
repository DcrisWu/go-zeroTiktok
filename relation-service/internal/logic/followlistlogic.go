package logic

import (
	"context"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/models/pack"
	"google.golang.org/grpc/status"

	"go-zeroTiktok/relation-service/internal/svc"
	"go-zeroTiktok/relation-service/pb/relation"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowListLogic) FollowList(in *relation.FollowListReq) (*relation.FollowListResp, error) {
	followingList, err := db.ListFollowing(l.ctx, l.svcCtx.DB, in.UserId)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	list, err := pack.FollowingList(l.ctx, l.svcCtx.DB, followingList, in.Uid)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	return &relation.FollowListResp{
		UserList: list,
	}, nil
}
