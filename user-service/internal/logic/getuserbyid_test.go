package logic

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-zeroTiktok/user-service/internal/config"
	"go-zeroTiktok/user-service/internal/svc"
	"go-zeroTiktok/user-service/pb/user"
	"testing"
)

func TestGetUserById(t *testing.T) {
	c := config.Config{
		DataSource: "root:Wu9121522521@@tcp(localhost:3306)/tiktok?parseTime=true",
		Argon2ID: &config.Argon2Params{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
	svcCtx := svc.NewServiceContext(c)
	logic := NewGetUserByIdLogic(context.Background(), svcCtx)
	resp, err := logic.GetUserById(&user.UserReq{
		Uid:    2,
		UserId: 1,
	})
	assert.NoError(t, err)
	fmt.Printf("%+v", resp)
}
