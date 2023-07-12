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

func TestLogin(t *testing.T) {
	c := config.Config{
		DataSource: "root:123456@tcp(localhost:23306)/tiktok?parseTime=true",
		Argon2ID: &config.Argon2Params{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
	svcCtx := svc.NewServiceContext(c)
	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &user.LoginReq{
		UserName: "dcris1",
		Password: "test_password",
	}
	register, err := logic.Login(req)
	assert.NoError(t, err)
	logx.Info(register)
}
