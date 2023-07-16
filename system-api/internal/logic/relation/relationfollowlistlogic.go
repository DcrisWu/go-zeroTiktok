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

type RelationFollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationFollowListLogic {
	return &RelationFollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationFollowListLogic) RelationFollowList(req *types.RelationFollowListReq) (*types.RelationFollowListResp, error) {
	uid := utils.GetUid(l.ctx)
	isExit, err := utils.IsJwtInRedis(l.ctx, l.svcCtx.Redis, uid)
	if err != nil || !isExit {
		return &types.RelationFollowListResp{
			StatusCode: http.StatusUnauthorized,
			StatusMsg:  "请先登录",
		}, nil
	}
	if uid == utils.UidNotFound || req.UserId == "0" {
		return &types.RelationFollowListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return &types.RelationFollowListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	list, err := l.svcCtx.RelationService.FollowList(l.ctx, &relation.FollowListReq{
		Uid:    uid,
		UserId: userId,
	})
	if err != nil {
		return &types.RelationFollowListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "获取关注列表失败",
		}, nil
	}
	followingList := make([]*types.User, 0)
	for _, u := range list.UserList {
		followingList = append(followingList, &types.User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      u.IsFollow,
		})
	}
	return &types.RelationFollowListResp{
		StatusCode: utils.SUCCESS,
		StatusMsg:  "获取关注列表成功",
		UserList:   followingList,
	}, nil
}
