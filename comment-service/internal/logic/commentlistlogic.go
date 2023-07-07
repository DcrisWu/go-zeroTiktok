package logic

import (
	"context"
	"google.golang.org/grpc/status"

	"go-zeroTiktok/comment-service/internal/svc"
	"go-zeroTiktok/comment-service/pb/comment"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentListLogic) CommentList(in *comment.CommentListReq) (*comment.CommentListResp, error) {
	comments, err := l.svcCtx.CommentModel.ListCommentByVideoId(l.ctx, in.VideoId)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	commentList := make([]*comment.CommentInfo, 0)
	for _, commentInfo := range comments {
		commentList = append(commentList, &comment.CommentInfo{
			CommentId: commentInfo.Id,
			UserId:    commentInfo.UserId,
			Content:   commentInfo.Content,
			CreatedAt: commentInfo.CreatedAt.Format("2006-01-02"),
		})
	}
	return &comment.CommentListResp{
		CommentList: commentList,
	}, nil
}
