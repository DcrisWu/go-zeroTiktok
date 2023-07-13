package logic

import (
	"context"
	"go-zeroTiktok/models/db"
	"google.golang.org/grpc/status"

	"go-zeroTiktok/relation-service/internal/svc"
	"go-zeroTiktok/relation-service/pb/relation"

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

func (l *ActionLogic) Action(in *relation.ActionReq) (*relation.ActionResp, error) {
	if in.ActionType == 1 {
		err := db.CreateRelation(l.ctx, l.svcCtx.DB, in.UserId, in.ToUserId)
		if err != nil {
			return nil, err
		}
		return &relation.ActionResp{}, nil
	}
	if in.ActionType == 2 {
		err := db.CancelRelation(l.ctx, l.svcCtx.DB, in.UserId, in.ToUserId)
		if err != nil {
			return nil, err
		}
		return &relation.ActionResp{}, nil
	}
	return nil, status.Error(1000, "参数错误")
}
