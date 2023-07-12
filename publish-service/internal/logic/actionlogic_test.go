package logic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/publish-service/internal/config"
	"go-zeroTiktok/publish-service/internal/svc"
	"go-zeroTiktok/publish-service/pb/publish"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
)

func NewServiceContext4Test() *svc.ServiceContext {
	c := config.Config{
		DataSource: "root:123456@tcp(localhost:23306)/tiktok?parseTime=true",
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

func TestAction(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewActionLogic(context.Background(), svcCtx)
	req := &publish.ActionReq{
		AuthorId: 1,
		PlayUrl:  "http://www.baidu.com",
		CoverUrl: "http://www.baidu.com",
		Title:    "test1",
	}
	action, err := logic.Action(req)
	assert.NoError(t, err)
	fmt.Printf("%+v", action)
}
