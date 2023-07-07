package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zeroTiktok/feed-service/internal/config"
	"go-zeroTiktok/model"
)

type ServiceContext struct {
	Config        config.Config
	VideoModel    model.VideoModel
	UserModel     model.UserModel
	FavoriteModel model.FavoriteModel
	RelationModel model.RelationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:        c,
		VideoModel:    model.NewVideoModel(sqlConn),
		UserModel:     model.NewUserModel(sqlConn),
		FavoriteModel: model.NewFavoriteModel(sqlConn),
		RelationModel: model.NewRelationModel(sqlConn),
	}
}
