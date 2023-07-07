package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zeroTiktok/model"
	"google.golang.org/grpc/status"

	"go-zeroTiktok/comment-service/internal/svc"
	"go-zeroTiktok/comment-service/pb/comment"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(in *comment.DeleteCommentReq) (*comment.DeleteCommentResp, error) {
	com, err := l.svcCtx.CommentModel.FindOneByUserIdVideoId(l.ctx, in.Uid, in.VideoId)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, status.Error(500, err.Error())
	}
	if errors.Is(err, model.ErrNotFound) {
		return nil, status.Error(500, "评论不存在")
	}
	// todo 事务：需要同时修改video的评论
	// 删除评论
	err = l.svcCtx.CommentModel.DeleteByCommentId(l.ctx, com.Id)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	err = l.svcCtx.VideoModel.IncrCommentCount(l.ctx, -1, com.VideoId)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	return &comment.DeleteCommentResp{
		IsDeleted: true,
	}, nil
}
