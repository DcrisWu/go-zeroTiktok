package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zeroTiktok/model"
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

	id, err := l.CreateUser(in, l.svcCtx.Config.Argon2ID)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		} else {
			return nil, status.Error(500, "注册失败")
		}
	}
	//注册成功直接登陆
	_, err = l.CheckUser(in)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		} else {
			return nil, status.Error(500, "注册失败")
		}
	}
	return &user.RegisterResp{
		Status: utils.SUCCESS,
		UserId: id,
	}, nil
}

func (l *RegisterLogic) CreateUser(in *user.RegisterReq, argon2Params *config.Argon2Params) (int64, error) {
	_, err := l.svcCtx.UserModel.FindOneByUserName(l.ctx, in.UserName)
	if err != nil && err != model.ErrNotFound {
		logx.Error(err)
		return 0, err
	}
	password, err := logic.GenerateFromPassword(in.Password, argon2Params)
	if err != nil {
		return 0, err
	}
	id, err := utils.NewBasicGenerator().GenerateId()
	if err != nil {
		return 0, err
	}
	_, err = l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Id:       id,
		UserName: in.UserName,
		Password: password,
	})
	if err != nil {
		logx.Error(err)
		return 0, err

	}
	return id, nil
}

func (l *RegisterLogic) CheckUser(in *user.RegisterReq) (int64, error) {
	u, err := l.svcCtx.UserModel.FindOneByUserName(l.ctx, in.UserName)
	if err != nil {
		if err == model.ErrNotFound {
			return 0, status.Error(400, "用户不存在")
		} else {
			return 0, status.Error(500, "注册失败")
		}
	}
	match, err := logic.ComparePasswordAndHash(in.Password, u.Password)
	if err != nil {
		return utils.UidNotFound, err
	}
	if !match {
		return 0, status.Error(400, "密码错误")
	}
	return u.Id, nil
}
