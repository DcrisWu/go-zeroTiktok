package logic

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/user-service/internal/config"
	"go-zeroTiktok/user-service/internal/svc"
	"go-zeroTiktok/user-service/pb/user"
	"testing"
)

func TestRegister(t *testing.T) {
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
	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &user.RegisterReq{
		UserName: "test_name1",
		Password: "test_password",
	}
	register, err := logic.Register(req)
	assert.NoError(t, err)
	logx.Info(register)
}
