package logic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-zeroTiktok/feed-service/internal/config"
	"go-zeroTiktok/feed-service/internal/svc"
	"go-zeroTiktok/feed-service/pb/feed"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
)

func TestFeedLogic(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewFeedLogic(context.Background(), svcCtx)
	in := &feed.FeedReq{
		Uid: utils.UidNotFound,
	}
	resp, err := logic.Feed(in)
	assert.NoError(t, err)
	fmt.Printf("resp:%v", resp)
}

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
