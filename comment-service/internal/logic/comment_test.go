package logic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-zeroTiktok/comment-service/internal/config"
	"go-zeroTiktok/comment-service/internal/svc"
	"go-zeroTiktok/comment-service/pb/comment"
	"go-zeroTiktok/models/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
)

func NewServiceContext4Test() *svc.ServiceContext {
	c := config.Config{
		DataSource: "root:Wu9121522521@@tcp(localhost:3306)/tiktok?parseTime=true",
	}
	database, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	database.AutoMigrate(&db.Comment{}, db.User{}, db.Relation{})
	return &svc.ServiceContext{
		Config: c,
		DB:     database,
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
