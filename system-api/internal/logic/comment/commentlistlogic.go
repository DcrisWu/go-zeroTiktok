package comment

import (
	"context"
	"go-zeroTiktok/comment-service/pb/comment"
	"go-zeroTiktok/utils"
	"net/http"
	"strconv"

	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (*types.CommentListResp, error) {
	uid := utils.GetUid(l.ctx)
	isExit, err := utils.IsJwtInRedis(l.ctx, l.svcCtx.Redis, uid)
	if err != nil || !isExit {
		return &types.CommentListResp{
			StatusCode: http.StatusUnauthorized,
			StatusMsg:  "请先登录",
		}, nil
	}
	if uid == utils.UidNotFound || req.VedioId == "" {
		return &types.CommentListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	vid, err := strconv.ParseInt(req.VedioId, 10, 64)
	if err != nil {
		return &types.CommentListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	list, err := l.svcCtx.CommentService.CommentList(l.ctx, &comment.CommentListReq{
		VideoId: vid,
	})
	if err != nil {
		return &types.CommentListResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "获取评论失败",
		}, nil
	}
	commentList := make([]*types.Comment, 0)
	for _, commentInfo := range list.CommentList {
		userInfo := &types.User{
			Id:            commentInfo.User.Id,
			Name:          commentInfo.User.Name,
			FollowCount:   commentInfo.User.FollowCount,
			FollowerCount: commentInfo.User.FollowerCount,
			IsFollow:      commentInfo.User.IsFollow,
		}
		commentList = append(commentList, &types.Comment{
			Id:         commentInfo.CommentId,
			UserInfo:   userInfo,
			Content:    commentInfo.Content,
			CreateDate: commentInfo.CreatedAt,
		})
	}
	return &types.CommentListResp{
		StatusCode:  utils.SUCCESS,
		StatusMsg:   "获取评论成功",
		CommentList: commentList,
	}, nil
}
