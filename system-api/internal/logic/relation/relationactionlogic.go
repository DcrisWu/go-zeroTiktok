package relation

import (
	"context"
	"go-zeroTiktok/relation-service/pb/relation"
	"go-zeroTiktok/utils"
	"net/http"
	"strconv"

	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationActionLogic {
	return &RelationActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationActionLogic) RelationAction(req *types.RelationActionReq) (*types.RelationActionResp, error) {
	uid := utils.GetUid(l.ctx)
	isExit, err := utils.IsJwtInRedis(l.ctx, l.svcCtx.Redis, uid)
	if err != nil || !isExit {
		return &types.RelationActionResp{
			StatusCode: http.StatusUnauthorized,
			StatusMsg:  "请先登录",
		}, nil
	}
	if uid == utils.UidNotFound || req.ToUserId == "0" || (req.ActionType != "1" && req.ActionType != "2") {
		return &types.RelationActionResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	toUserId, err := strconv.ParseInt(req.ActionType, 10, 64)
	if err != nil {
		return &types.RelationActionResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	_, err = l.svcCtx.RelationService.Action(l.ctx, &relation.ActionReq{
		UserId:   uid,
		ToUserId: toUserId,
	})
	if err != nil {
		return &types.RelationActionResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "关注失败",
		}, nil
	}
	return &types.RelationActionResp{
		StatusCode: utils.SUCCESS,
		StatusMsg:  "关注成功",
	}, nil
}
