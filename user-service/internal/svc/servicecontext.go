package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/user-service/internal/config"
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
	//sqlConn := sqlx.NewMysql(c.DataSource)

	database, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	database.AutoMigrate(&db.Comment{}, db.User{}, db.Relation{}, db.Video{})
	return &ServiceContext{
		Config: c,
		DB:     database,
		Redis:  redis.MustNewRedis(c.RedisCfg),
	}
}
