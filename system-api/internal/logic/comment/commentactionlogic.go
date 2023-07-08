package comment

import (
	"context"
	"go-zeroTiktok/comment-service/pb/comment"
	"go-zeroTiktok/user-service/pb/user"
	"go-zeroTiktok/utils"
	"strconv"

	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {
	return &CommentActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentActionLogic) CommentAction(req *types.CommentActionReq) (*types.CommentActionResp, error) {
	uid := utils.GetUid(l.ctx)
	if (req.ActionType != "1" && req.ActionType != "2") || uid == utils.PayLoadNotFound ||
		uid == utils.UidNotFound || req.VideoId == "" {
		return &types.CommentActionResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "参数错误",
		}, nil
	}
	if req.ActionType == "1" {
		if req.CommentText == "" {
			return &types.CommentActionResp{
				StatusCode: utils.FAILED,
				StatusMsg:  "参数错误",
			}, nil
		}
		parseInt, err := strconv.ParseInt(req.VideoId, 10, 64)
		if err != nil {
			return &types.CommentActionResp{
				StatusCode: utils.FAILED,
				StatusMsg:  "评论失败",
			}, nil
		}
		commentReq := &comment.CreateCommentReq{
			Uid:     uid,
			VideoId: parseInt,
			Content: req.CommentText,
		}
		resp, err := l.svcCtx.CommentService.CreateComment(l.ctx, commentReq)
		if err != nil {
			return &types.CommentActionResp{
				StatusCode: utils.FAILED,
				StatusMsg:  "评论失败",
			}, nil
		}
		userResp, err := l.svcCtx.UserService.GetUserById(l.ctx, &user.UserReq{
			Uid:    uid,
			UserId: uid,
		})
		userInfo := &types.User{}
		if err != nil {
			logx.Error("获取用户信息错误, ", err.Error())
		} else {
			userInfo.Id = userResp.User.Id
			userInfo.Name = userResp.User.Name
			userInfo.FollowCount = userResp.User.FollowCount
			userInfo.FollowerCount = userResp.User.FollowerCount
			userInfo.IsFollow = userResp.User.IsFollow
		}
		return &types.CommentActionResp{
			StatusCode: utils.SUCCESS,
			StatusMsg:  "评论成功",
			CommentObj: &types.Comment{
				UserInfo:   userInfo,
				Content:    req.CommentText,
				CreateDate: resp.CreatedAt,
			},
		}, nil
	}
	if req.ActionType == "2" {
		if req.CommentId == "" {
			return &types.CommentActionResp{
				StatusCode: utils.FAILED,
				StatusMsg:  "参数错误",
			}, nil
		}
		parseInt, err := strconv.ParseInt(req.VideoId, 10, 64)
		if err != nil {
			return &types.CommentActionResp{
				StatusCode: utils.FAILED,
				StatusMsg:  "评论删除失败",
			}, nil
		}
		commentId, err := strconv.ParseInt(req.CommentId, 10, 64)
		if err != nil {
			return &types.CommentActionResp{
				StatusCode: utils.FAILED,
				StatusMsg:  "评论删除失败",
			}, nil
		}
		commentReq := &comment.DeleteCommentReq{
			CommentId: commentId,
			VideoId:   parseInt,
		}
		resp, err := l.svcCtx.CommentService.DeleteComment(l.ctx, commentReq)
		if err != nil || !resp.IsDeleted {
			return &types.CommentActionResp{
				StatusCode: utils.FAILED,
				StatusMsg:  "评论删除失败",
			}, nil
		}
		return &types.CommentActionResp{
			StatusCode: utils.SUCCESS,
			StatusMsg:  "评论删除成功",
		}, nil
	}
	return &types.CommentActionResp{
		StatusCode: utils.FAILED,
		StatusMsg:  "参数错误",
	}, nil
}
