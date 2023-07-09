package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zeroTiktok/favorite-service/internal/config"
	"go-zeroTiktok/models/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Redis
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
		Config: c,
		DB:     database,
		Redis:  redis.MustNewRedis(c.RedisCfg),
	}
}
