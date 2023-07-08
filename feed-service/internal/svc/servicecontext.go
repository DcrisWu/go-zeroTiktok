package svc

import (
	"go-zeroTiktok/feed-service/internal/config"
	"go-zeroTiktok/models/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config config.Config
	//VideoModel    model.VideoModel
	//UserModel     model.UserModel
	//FavoriteModel model.FavoriteModel
	//RelationModel model.RelationModel
	DB *gorm.DB
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
	database.AutoMigrate(&db.Video{}, db.User{}, db.Relation{})
	return &ServiceContext{
		Config: c,
		DB:     database,
	}
}
