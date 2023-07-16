package logic

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-zeroTiktok/relation-service/pb/relation"
	"testing"
)

func TestFollowerList(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewFollowerListLogic(context.Background(), svcCtx)
	list, err := logic.FollowerList(&relation.FollowerListReq{
		Uid:    1,
		UserId: 2,
	})
	assert.NoError(t, err)
	t.Log(list)
}
