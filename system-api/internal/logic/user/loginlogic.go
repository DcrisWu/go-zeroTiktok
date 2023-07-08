package user

import (
	"context"
	"go-zeroTiktok/user-service/pb/user"
	"go-zeroTiktok/utils"
	"google.golang.org/grpc/status"
	"strconv"
	"time"

	"go-zeroTiktok/system-api/internal/svc"
	"go-zeroTiktok/system-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	if req.UserName == "" || req.Password == "" {
		return &types.LoginResp{
			StatusCode: utils.FAILED,
			StatusMsg:  "用户名或密码不能为空",
		}, nil
	}
	resp, err := l.svcCtx.UserService.Login(l.ctx, &user.LoginReq{
		UserName: req.UserName,
		Password: req.Password,
	})
	if err != nil {
		if s, ok := status.FromError(err); !ok {
			return &types.LoginResp{
				StatusCode: utils.FAILED,
				StatusMsg:  s.Message(),
			}, nil
		} else {
			return nil, err
		}
	}
	payload := make(map[string]interface{})
	// js 的 number 类型最大值为 2^53 - 1，超过这个值会丢失精度，所以这里需要转成字符串
	payload["uid"] = strconv.FormatInt(resp.UserId, 10)
	token, err := utils.GenerateJwt(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, payload)
	return &types.LoginResp{
		StatusCode: utils.SUCCESS,
		StatusMsg:  "登陆成功",
		UserId:     resp.UserId,
		Token:      token,
	}, nil
}
