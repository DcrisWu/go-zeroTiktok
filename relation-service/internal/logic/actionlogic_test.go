package logic

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/relation-service/internal/config"
	"go-zeroTiktok/relation-service/internal/svc"
	"go-zeroTiktok/relation-service/pb/relation"
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
	database.AutoMigrate(&db.Video{}, db.User{}, db.Comment{}, &db.Relation{})

	return &svc.ServiceContext{
		Config: c,
		DB:     database,
	}
}

func TestCreateActionLogic(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewActionLogic(context.Background(), svcCtx)
	_, err := logic.Action(&relation.ActionReq{
		UserId:     3,
		ToUserId:   2,
		ActionType: 1,
	})
	assert.NoError(t, err)
}

func TestCancelActionLogic(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewActionLogic(context.Background(), svcCtx)
	_, err := logic.Action(&relation.ActionReq{
		UserId:     2,
		ToUserId:   1,
		ActionType: 2,
	})
	assert.NoError(t, err)
}
