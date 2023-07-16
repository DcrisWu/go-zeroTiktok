package publish

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/publish-service/pb/publish"
	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"
	"go-zeroTiktok/utils"
	"net/http"
)

type ActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionLogic {
	return &ActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActionLogic) Action(req *types.PublishActionReq) (*types.PublishActionResp, error) {
	uid := utils.GetUid(l.ctx)
	isExit, err := utils.IsJwtInRedis(l.ctx, l.svcCtx.Redis, uid)
	if err != nil || !isExit {
		return &types.PublishActionResp{
			StatusCode: http.StatusUnauthorized,
			StatusMsg:  "请先登录",
		}, nil
	}
	if req.Title == "" || req.Data.PlayUrl == "" || req.Data.CoverUrl == "" {
		return &types.PublishActionResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	action, err := l.svcCtx.PublishService.Action(l.ctx, &publish.ActionReq{
		AuthorId: uid,
		PlayUrl:  req.Data.PlayUrl,
		CoverUrl: req.Data.CoverUrl,
		Title:    req.Title,
	})
	if err != nil || !action.IsSuccess {
		return &types.PublishActionResp{
			StatusCode: utils.FAILED,
			StatusMsg:  err.Error(),
		}, nil
	}

	return &types.PublishActionResp{
		StatusCode: utils.SUCCESS,
		StatusMsg:  "发布成功",
	}, nil
}
