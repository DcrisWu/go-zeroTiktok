package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zeroTiktok/favorite-service/internal/config"
	"go-zeroTiktok/favorite-service/internal/logic/favoritemq"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	Redis      *redis.Redis
	FavoriteMq *utils.RabbitMq
}

func NewServiceContext(c config.Config) *ServiceContext {
	database, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	database.AutoMigrate(&db.Video{}, db.User{}, db.Comment{}, &db.Relation{})

	return &ServiceContext{
		Config:     c,
		DB:         database,
		Redis:      redis.MustNewRedis(c.RedisCfg),
		FavoriteMq: InitFavoriteMq(c.MqUrl, database),
	}
}

func InitFavoriteMq(mqUrl string, DB *gorm.DB) *utils.RabbitMq {
	utils.InitRabbitMQ(mqUrl)
	mq := utils.NewRabbitMq("favorite_queue", "favorite_exchange", "favorite")
	// 启动mq消费者
	go favoritemq.FavoriteConsumer(mq, DB)
	return mq
}
