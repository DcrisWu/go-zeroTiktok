package logic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zeroTiktok/feed-service/internal/config"
	"go-zeroTiktok/feed-service/internal/svc"
	"go-zeroTiktok/feed-service/pb/feed"
	"go-zeroTiktok/model"
	"go-zeroTiktok/utils"
	"testing"
)

func TestFeedLogic(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewFeedLogic(context.Background(), svcCtx)
	in := &feed.FeedReq{
		Uid: utils.UidNotFound,
	}
	resp, err := logic.Feed(in)
	assert.NoError(t, err)
	fmt.Printf("resp:%v", resp)
}

func NewServiceContext4Test() *svc.ServiceContext {
	c := config.Config{
		DataSource: "root:Wu9121522521@@tcp(localhost:3306)/tiktok?parseTime=true",
	}
	return &svc.ServiceContext{
		Config:        c,
		VideoModel:    model.NewVideoModel(sqlx.NewMysql(c.DataSource)),
		UserModel:     model.NewUserModel(sqlx.NewMysql(c.DataSource)),
		FavoriteModel: model.NewFavoriteModel(sqlx.NewMysql(c.DataSource)),
		RelationModel: model.NewRelationModel(sqlx.NewMysql(c.DataSource)),
	}

}
