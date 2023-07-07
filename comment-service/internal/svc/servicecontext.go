package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zeroTiktok/comment-service/internal/config"
	"go-zeroTiktok/model"
)

type ServiceContext struct {
	Config       config.Config
	MysqlConn    sqlx.SqlConn
	CommentModel model.CommentModel
	VideoModel   model.VideoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:       c,
		MysqlConn:    sqlConn,
		CommentModel: model.NewCommentModel(sqlConn),
		VideoModel:   model.NewVideoModel(sqlConn),
	}
}
