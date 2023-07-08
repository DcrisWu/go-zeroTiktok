package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zeroTiktok/models/db"
	logic "go-zeroTiktok/user-service/internal/logic/userutils"
	"go-zeroTiktok/user-service/internal/svc"
	"go-zeroTiktok/user-service/pb/user"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	uid, err := l.CheckUser(in)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		} else {
			return nil, status.Error(500, "登陆失败")
		}
	}
	return &user.LoginResp{
		UserId: uid,
	}, nil
}

func (l *LoginLogic) CheckUser(in *user.LoginReq) (int64, error) {
	users, err := db.QueryUser(l.ctx, l.svcCtx.DB, in.UserName)
	if err != nil {
		logx.Error(err)
		return 0, errors.Errorf("数据库查询失败")
	}
	if len(users) == 0 {
		return 0, status.Error(400, "用户不存在")
	}
	u := users[0]
	match, err := logic.ComparePasswordAndHash(in.Password, u.Password)
	if err != nil {
		return 0, err
	}
	if !match {
		return 0, status.Error(400, "密码错误")
	}
	return int64(u.ID), nil
}
