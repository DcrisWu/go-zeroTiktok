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

type RelationFollowerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationFollowerListLogic {
	return &RelationFollowerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationFollowerListLogic) RelationFollowerList(req *types.RelationFollowerListReq) (*types.RelationFollowerListResp, error) {
	uid := utils.GetUid(l.ctx)
	isExit, err := utils.IsJwtInRedis(l.ctx, l.svcCtx.Redis, uid)
	if err != nil || !isExit {
		return &types.RelationFollowerListResp{
			StatusCode: http.StatusUnauthorized,
			StatusMsg:  "请先登录",
		}, nil
	}
	if uid == utils.UidNotFound || req.UserId == "0" {
		return &types.RelationFollowerListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return &types.RelationFollowerListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	list, err := l.svcCtx.RelationService.FollowerList(l.ctx, &relation.FollowerListReq{
		Uid:    uid,
		UserId: userId,
	})
	if err != nil {
		return &types.RelationFollowerListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "获取关注列表失败",
		}, nil
	}
	followerList := make([]*types.User, 0)
	for _, u := range list.UserList {
		followerList = append(followerList, &types.User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
		})
	}
	return &types.RelationFollowerListResp{
		StatusCode: utils.SUCCESS,
		StatusMsg:  "获取关注列表成功",
		UserList:   followerList,
	}, nil
}
