package logic

import (
	"context"
	"encoding/json"
	"go-zeroTiktok/relation-service/internal/logic/relationmq"
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
		msg, err := json.Marshal(in)
		if err != nil {
			return nil, status.Error(1000, "json.Marshal error:"+err.Error())
		}
		relationmq.RelationActionMqSend(l.svcCtx.RelationMq, msg)
		addRedisFollowList(l.ctx, l.svcCtx.DB, l.svcCtx.Redis, in.UserId, in.ToUserId)
		addRedisFollowerList(l.ctx, l.svcCtx.DB, l.svcCtx.Redis, in.UserId, in.ToUserId)
	}
	if in.ActionType == 2 {
		// 取消关注
		msg, err := json.Marshal(in)
		if err != nil {
			return nil, status.Error(1000, "json.Marshal error:"+err.Error())
		}
		relationmq.RelationActionMqSend(l.svcCtx.RelationMq, msg)
		rmRedisFollowList(l.ctx, l.svcCtx.DB, l.svcCtx.Redis, in.UserId, in.ToUserId)
		rmRedisFollowerList(l.ctx, l.svcCtx.DB, l.svcCtx.Redis, in.UserId, in.ToUserId)
	}
	return nil, status.Error(1000, "关注功能暂未开放")
}
