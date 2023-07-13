package logic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zeroTiktok/favorite-service/internal/config"
	"go-zeroTiktok/favorite-service/internal/logic/favoritemq"
	"go-zeroTiktok/favorite-service/internal/svc"
	"go-zeroTiktok/favorite-service/pb/favorite"
	"go-zeroTiktok/models/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
	"time"
)

func NewServiceContext4Test() *svc.ServiceContext {
	c := config.Config{
		DataSource: "root:123456@tcp(localhost:23306)/tiktok?parseTime=true",
		RedisCfg: redis.RedisConf{
			Host: "localhost:26379",
			Type: "single",
			Pass: "abcd",
		},
		MqUrl: "amqp://admin:123456@localhost:5672/",
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
		Config:     c,
		DB:         database,
		Redis:      redis.MustNewRedis(c.RedisCfg),
		FavoriteMq: svc.InitFavoriteMq(c.MqUrl, database),
	}
}

func TestCreateFavorite(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewFavoriteActionLogic(context.Background(), svcCtx)
	req := &favorite.FavoriteActionReq{
		ActionType: 1,
		UserId:     1,
		VideoId:    1,
	}
	_, err := logic.FavoriteAction(req)
	assert.NoError(t, err)
	go favoritemq.FavoriteConsumer(svcCtx.FavoriteMq, svcCtx.DB)
	time.Sleep(time.Second * 5)
}

func TestCancelFavorite(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewFavoriteActionLogic(context.Background(), svcCtx)
	req := &favorite.FavoriteActionReq{
		ActionType: 2,
		UserId:     1,
		VideoId:    1,
	}
	_, err := logic.FavoriteAction(req)
	assert.NoError(t, err)
	go favoritemq.FavoriteConsumer(svcCtx.FavoriteMq, svcCtx.DB)
	time.Sleep(time.Second * 5)
}

func TestFavoriteList(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewGetFavoriteListLogic(context.Background(), svcCtx)
	req := &favorite.FavoriteListReq{
		UserId: 2,
	}
	list, err := logic.GetFavoriteList(req)
	assert.NoError(t, err)
	fmt.Printf("%+v\n", list)
}
