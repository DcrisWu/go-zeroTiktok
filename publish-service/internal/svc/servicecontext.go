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
	sqlConn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:        c,
		VideoModel:    model.NewVideoModel(sqlConn),
		UserModel:     model.NewUserModel(sqlConn),
		RelationModel: model.NewRelationModel(sqlConn),
		FavoriteModel: model.NewFavoriteModel(sqlConn),
	}
}
