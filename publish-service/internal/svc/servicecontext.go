package svc

import (
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/publish-service/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config config.Config
	//VideoModel    model.VideoModel
	//UserModel     model.UserModel
	//RelationModel model.RelationModel
	//FavoriteModel model.FavoriteModel
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
	database.AutoMigrate(db.User{}, db.Relation{})
	return &ServiceContext{
		Config: c,
		DB:     database,
	}
}
