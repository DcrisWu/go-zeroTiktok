package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/user-service/internal/config"
	"go-zeroTiktok/user-service/internal/logic/userutils"
	"go-zeroTiktok/user-service/internal/svc"
	"go-zeroTiktok/user-service/pb/user"
	"go-zeroTiktok/utils"
	"google.golang.org/grpc/status"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {

	err := l.CreateUser(in, l.svcCtx.Config.Argon2ID)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		} else {
			return nil, status.Error(500, "注册失败")
		}
	}
	//注册成功直接登陆
	uid, err := l.CheckUser(in)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		} else {
			return nil, status.Error(500, "注册失败")
		}
	}
	return &user.RegisterResp{
		Status: utils.SUCCESS,
		UserId: uid,
	}, nil
}

func (l *RegisterLogic) CreateUser(in *user.RegisterReq, argon2Params *config.Argon2Params) error {
	//u, err := l.svcCtx.UserModel.FindOneByUserName(l.ctx, in.UserName)
	users, err := db.QueryUser(l.ctx, l.svcCtx.DB, in.UserName)
	if err != nil {
		return err
	}
	if len(users) != 0 {
		return errors.Errorf("用户已存在")
	}

	password, err := logic.GenerateFromPassword(in.Password, argon2Params)
	if err != nil {
		return err
	}
	return db.CreateUser(l.ctx, l.svcCtx.DB, []*db.User{{
		UserName: in.UserName,
		Password: password,
	}})
}

func (l *RegisterLogic) CheckUser(in *user.RegisterReq) (int64, error) {

	users, err := db.QueryUser(l.ctx, l.svcCtx.DB, in.UserName)
	if err != nil {
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
