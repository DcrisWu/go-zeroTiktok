package svc

import (
	"go-zeroTiktok/comment-service/internal/config"
	"go-zeroTiktok/models/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config config.Config
	//MysqlConn    sqlx.SqlConn
	//CommentModel model.CommentModel
	//VideoModel   model.VideoModel
	DB *gorm.DB
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
	}
}
