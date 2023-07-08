package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/models/pack"
	"go-zeroTiktok/user-service/internal/svc"
	"go-zeroTiktok/user-service/pb/user"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByIdLogic) GetUserById(in *user.UserReq) (*user.UserResp, error) {
	if in.Uid == 0 || in.UserId == 0 {
		return nil, status.Error(400, "参数缺失")
	}
	//one, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	modelUser, err := db.GetUserById(l.ctx, l.svcCtx.DB, in.UserId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(500, err.Error())
	}
	// 获取 uid 是否关注了 userId，将 db.User 封装为 user.User
	userInfo, err := pack.User(l.ctx, l.svcCtx.DB, modelUser, in.Uid)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	return &user.UserResp{
		User: userInfo,
	}, nil
}
