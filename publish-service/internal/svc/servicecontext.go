package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zeroTiktok/model"
	"go-zeroTiktok/publish-service/internal/config"
)

type ServiceContext struct {
	Config        config.Config
	VideoModel    model.VideoModel
	UserModel     model.UserModel
	RelationModel model.RelationModel
	FavoriteModel model.FavoriteModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		VideoModel:    model.NewVideoModel(sqlx.NewMysql(c.DataSource)),
		UserModel:     model.NewUserModel(sqlx.NewMysql(c.DataSource)),
		RelationModel: model.NewRelationModel(sqlx.NewMysql(c.DataSource)),
		FavoriteModel: model.NewFavoriteModel(sqlx.NewMysql(c.DataSource)),
	}
}
