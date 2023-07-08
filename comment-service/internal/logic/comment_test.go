package logic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zeroTiktok/comment-service/internal/config"
	"go-zeroTiktok/comment-service/internal/svc"
	"go-zeroTiktok/comment-service/pb/comment"
	"go-zeroTiktok/model"
	"testing"
)

func NewServiceContext4Test() *svc.ServiceContext {
	c := &config.Config{
		DataSource: "root:Wu9121522521@@tcp(localhost:3306)/tiktok?parseTime=true",
	}
	sqlConn := sqlx.NewMysql(c.DataSource)
	return &svc.ServiceContext{
		MysqlConn:    sqlConn,
		CommentModel: model.NewCommentModel(sqlConn),
		VideoModel:   model.NewVideoModel(sqlConn),
	}
}

func TestCreateCommentLogic_CreateComment(t *testing.T) {
	req := &comment.CreateCommentReq{
		Uid:     6828874698114169068,
		VideoId: 3321574323690790359,
		Content: "fucking high",
	}
	svcCtx := NewServiceContext4Test()
	logic := NewCreateCommentLogic(context.Background(), svcCtx)
	createComment, err := logic.CreateComment(req)
	assert.NoError(t, err)
	fmt.Printf("%+v", createComment)
}
