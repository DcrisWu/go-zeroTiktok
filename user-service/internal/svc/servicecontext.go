package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zeroTiktok/model"
	"go-zeroTiktok/user-service/internal/config"
)

type ServiceContext struct {
	Config        config.Config
	UserModel     model.UserModel
	RelationModel model.RelationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:        c,
		UserModel:     model.NewUserModel(sqlConn),
		RelationModel: model.NewRelationModel(sqlConn),
	}
}
