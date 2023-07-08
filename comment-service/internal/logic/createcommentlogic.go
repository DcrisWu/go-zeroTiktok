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
		UserId:  in.Uid,
		VideoId: in.VideoId,
		Content: in.Content,
	})
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	return &comment.CreateCommentResp{
		CreatedAt: time.Now().Format("2023-07-08"),
	}, nil

	//id, err := utils.NewBasicGenerator().GenerateCommentId()
	//if err != nil {
	//	return nil, status.Error(500, "id生成失败")
	//}
	//// 事务
	//err = l.svcCtx.MysqlConn.Transact(func(session sqlx.Session) error {
	//	m := &model.Comment{
	//		//Id:      id,
	//		UserId:  in.Uid, // 评论那个人的id
	//		VideoId: in.VideoId,
	//		Content: in.Content,
	//	}
	//	// 新建评论
	//	err = l.svcCtx.CommentModel.InsertComment(l.ctx, m)
	//	if err != nil {
	//		return err
	//	}
	//	// 视频的评论数+1
	//	err = l.svcCtx.VideoModel.IncrCommentCount(l.ctx, 1, in.VideoId)
	//	if err != nil {
	//		return err
	//	}
	//	return nil
	//})
	//if err != nil {
	//	return nil, status.Error(500, err.Error())
	//}
	//
	////m := &model.Comment{
	////	Id:      id,
	////	UserId:  in.Uid, // 评论那个人的id
	////	VideoId: in.VideoId,
	////	Content: in.Content,
	////}
	////// 新建评论
	////_, err = l.svcCtx.CommentModel.Insert(l.ctx, m)
	////if err != nil {
	////	return nil, err
	////}
	////// 视频的评论数+1
	////err = l.svcCtx.VideoModel.IncrCommentCount(l.ctx, 1, in.VideoId)
	////if err != nil {
	////	return nil, err
	////}
	//return &comment.CreateCommentResp{
	//	CommentId: id,
	//	CreatedAt: time.Now().Format("2023-07-07"),
	//}, nil
}
