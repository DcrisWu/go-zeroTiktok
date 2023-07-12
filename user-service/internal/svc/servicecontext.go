package svc

import (
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/user-service/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config config.Config
	//UserModel     model.UserModel
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
	database.AutoMigrate(&db.Comment{}, db.User{}, db.Relation{}, db.Video{})
	return &ServiceContext{
		Config: c,
		//UserModel:     model.NewUserModel(sqlConn),
		//RelationModel: model.NewRelationModel(sqlConn),
		DB: database,
	}
}
