package logic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zeroTiktok/model"
	"go-zeroTiktok/publish-service/internal/config"
	"go-zeroTiktok/publish-service/internal/svc"
	"go-zeroTiktok/publish-service/pb/publish"
	"testing"
)

func NewServiceContext4Test() *svc.ServiceContext {
	c := config.Config{
		DataSource: "root:Wu9121522521@@tcp(localhost:3306)/tiktok?parseTime=true",
	}
	return &svc.ServiceContext{
		Config:        c,
		VideoModel:    model.NewVideoModel(sqlx.NewMysql(c.DataSource)),
		UserModel:     model.NewUserModel(sqlx.NewMysql(c.DataSource)),
		RelationModel: model.NewRelationModel(sqlx.NewMysql(c.DataSource)),
		FavoriteModel: model.NewFavoriteModel(sqlx.NewMysql(c.DataSource)),
	}
}

func TestAction(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewActionLogic(context.Background(), svcCtx)
	req := &publish.ActionReq{
		AuthorId: 6999740003925172302,
		PlayUrl:  "http://www.bilibili.com",
		CoverUrl: "http://www.bilibili.com",
		Title:    "test2",
	}
	action, err := logic.Action(req)
	assert.NoError(t, err)
	fmt.Printf("%+v", action)
}
