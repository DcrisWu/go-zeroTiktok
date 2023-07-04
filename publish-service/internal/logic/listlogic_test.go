package logic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-zeroTiktok/publish-service/pb/publish"
	"testing"
)

func TestListLogic(t *testing.T) {
	svcCtx := NewServiceContext4Test()
	logic := NewListLogic(context.Background(), svcCtx)
	// 6999740003925172302
	req := &publish.ListReq{
		AuthorId: 10,
	}
	list, err := logic.List(req)
	assert.NoError(t, err)
	fmt.Printf("%+v", list)
}
