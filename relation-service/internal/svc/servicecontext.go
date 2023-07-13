package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/relation-service/internal/config"
	"go-zeroTiktok/relation-service/internal/logic/relationmq"
	"go-zeroTiktok/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	Redis      *redis.Redis
	RelationMq *utils.RabbitMq
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
		RelationMq: InitRelationMq(c.MqUrl, database),
	}
}

func InitRelationMq(mqUrl string, DB *gorm.DB) *utils.RabbitMq {
	utils.InitRabbitMQ(mqUrl)
	mq := utils.NewRabbitMq("relation_queue", "relation_exchange", "relation")
	go relationmq.RelationConsumer(mq, DB)
	return mq
}
