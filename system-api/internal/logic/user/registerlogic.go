package user

import (
	"context"
	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"
	"go-zeroTiktok/user-service/pb/user"
	"go-zeroTiktok/utils"
	"google.golang.org/grpc/status"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (*types.RegisterResp, error) {
	if req.UserName == "" || req.Password == "" {
		return &types.RegisterResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "用户名或密码不能为空",
		}, nil
	}
	resp, err := l.svcCtx.UserService.Register(l.ctx, &user.RegisterReq{
		UserName: req.UserName,
		Password: req.Password,
	})
	if err != nil {
		if s, ok := status.FromError(err); !ok {
			return &types.RegisterResp{
				StatusCode: utils.FAILED,
				StatusMsg:  s.Message(),
			}, nil
		} else {
			return nil, err
		}
	}
	payload := make(map[string]interface{})
	payload["uid"] = resp.UserId
	token, err := utils.GenerateToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix()+l.svcCtx.Config.Auth.AccessExpire, l.svcCtx.Config.Auth.AccessExpire, payload)
	return &types.RegisterResp{
		StatusCode: utils.SUCCESS,
		StatusMsg:  "注册成功",
		UserId:     resp.UserId,
		Token:      token,
	}, nil
}
