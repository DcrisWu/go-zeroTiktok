package logic

import (
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/relation-service/internal/config"
	"go-zeroTiktok/relation-service/internal/svc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
