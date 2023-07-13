package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/comment-service/internal/svc"
	"go-zeroTiktok/comment-service/pb/comment"
	"go-zeroTiktok/models/db"
	"google.golang.org/grpc/status"
	"time"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCommentLogic) CreateComment(in *comment.CreateCommentReq) (*comment.CreateCommentResp, error) {
	err := db.NewComment(l.ctx, l.svcCtx.DB, &db.Comment{
		UserId:  int(in.Uid),
		VideoId: int(in.VideoId),
		Content: in.Content,
	})
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	return &comment.CreateCommentResp{
		CreatedAt: time.Now().Format("2023-07-08"),
	}, nil
}
