package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/comment-service/internal/svc"
	"go-zeroTiktok/comment-service/pb/comment"
	"go-zeroTiktok/model"
	"go-zeroTiktok/utils"
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
	id, err := utils.NewBasicGenerator().GenerateCommentId()
	if err != nil {
		return nil, status.Error(500, "id生成失败")
	}
	// todo 事务
	m := &model.Comment{
		Id:      id,
		UserId:  in.Uid, // 评论那个人的id
		VideoId: in.VideoId,
		Content: in.Content,
	}
	// 新建评论
	_, err = l.svcCtx.CommentModel.Insert(l.ctx, m)
	if err != nil {
		return nil, err
	}
	// 视频的评论数+1
	err = l.svcCtx.VideoModel.IncrCommentCount(l.ctx, 1, in.VideoId)
	if err != nil {
		return nil, err
	}
	return &comment.CreateCommentResp{
		CommentId: id,
		CreatedAt: time.Now().Format("2023-07-07"),
	}, nil
}
