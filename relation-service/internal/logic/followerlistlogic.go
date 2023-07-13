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

type FollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowerListLogic {
	return &FollowerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowerListLogic) FollowerList(in *relation.FollowerListReq) (*relation.FollowerListResp, error) {
	followerList, err := db.ListFollower(l.ctx, l.svcCtx.DB, in.UserId)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	list, err := pack.FollowerList(l.ctx, l.svcCtx.DB, followerList, in.Uid)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	return &relation.FollowerListResp{
		UserList: list,
	}, nil
}
