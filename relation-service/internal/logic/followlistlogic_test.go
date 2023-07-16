package logic

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-zeroTiktok/relation-service/pb/relation"
	"testing"
)

func TestFollowListLogic(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewFollowListLogic(context.Background(), svcCtx)
	list, err := logic.FollowList(&relation.FollowListReq{
		Uid:    1,
		UserId: 3,
	})
	assert.NoError(t, err)
	t.Log(list)
}
