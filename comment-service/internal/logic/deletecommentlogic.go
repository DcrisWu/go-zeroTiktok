package logic

import (
	"context"
	"go-zeroTiktok/comment-service/internal/svc"
	"go-zeroTiktok/comment-service/pb/comment"
	"go-zeroTiktok/models/db"
	"google.golang.org/grpc/status"

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
	err := db.DeleteComment(l.ctx, l.svcCtx.DB, in.CommentId, in.VideoId)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	return &comment.DeleteCommentResp{
		IsDeleted: true,
	}, nil

	//com, err := l.svcCtx.CommentModel.FindOneByUserIdVideoId(l.ctx, in.Uid, in.VideoId)
	//if err != nil && !errors.Is(err, model.ErrNotFound) {
	//	return nil, status.Error(500, err.Error())
	//}
	//if errors.Is(err, model.ErrNotFound) {
	//	return nil, status.Error(500, "评论不存在")
	//}
	//err = l.svcCtx.MysqlConn.Transact(func(session sqlx.Session) error {
	//	// 删除评论
	//	err = l.svcCtx.CommentModel.DeleteByCommentId(l.ctx, com.Id)
	//	if err != nil {
	//		return err
	//	}
	//	err = l.svcCtx.VideoModel.IncrCommentCount(l.ctx, -1, com.VideoId)
	//	return err
	//})
	//if err != nil {
	//	return nil, status.Error(500, err.Error())
	//}
	//return &comment.DeleteCommentResp{
	//	IsDeleted: true,
	//}, nil
}
