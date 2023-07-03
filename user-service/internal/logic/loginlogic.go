package logic

import (
	"context"
	logic "go-zeroTiktok/user-service/internal/logic/userutils"
	"go-zeroTiktok/user-service/internal/svc"
	"go-zeroTiktok/user-service/pb/user"
	"go-zeroTiktok/utils"
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
		Status: utils.SUCCESS,
		UserId: uid,
	}, nil
}

func (l *LoginLogic) CheckUser(in *user.LoginReq) (int64, error) {
	u, err := l.svcCtx.UserModel.FindOneByUserName(l.ctx, in.UserName)
	if err != nil {
		logx.Error(err)
		return 0, err
	}
	if u == nil {
		return 0, status.Error(400, "用户不存在")
	}
	match, err := logic.ComparePasswordAndHash(in.Password, u.Password)
	if err != nil {
		return 0, err
	}
	if !match {
		return 0, status.Error(400, "密码错误")
	}
	return u.Id, nil
}
