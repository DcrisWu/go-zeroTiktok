package logic

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/relation-service/internal/config"
	"go-zeroTiktok/relation-service/internal/svc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	redis := svcCtx.Redis
	for i := 1; i < 10; i++ {
		redis.Sadd("test", i)
		redis.Expire("test", 10)
		time.Sleep(5 * time.Second)
	}

}

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
		RelationMq: svc.InitRelationMq(c.MqUrl, database),
	}
}
