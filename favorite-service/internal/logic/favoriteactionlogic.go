package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go-zeroTiktok/favorite-service/internal/logic/favoritemq"
	"google.golang.org/grpc/status"
	"strconv"

	"go-zeroTiktok/favorite-service/internal/svc"
	"go-zeroTiktok/favorite-service/pb/favorite"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var prefix = "video-user-relation:"

func (l *FavoriteActionLogic) FavoriteAction(in *favorite.FavoriteActionReq) (*favorite.FavoriteActionResp, error) {
	if in.ActionType == 1 {
		err := l.SetFavoriteToRedis(in.UserId, in.VideoId)
		if err != nil {
			logx.Error(err)
			return nil, status.Error(500, err.Error())
		}

		// 发送消息给MQ
		msg, err := json.Marshal(in)
		if err != nil {
			logx.Error(err)
			return nil, status.Error(500, err.Error())
		}
		favoritemq.FavoriteActionSend(l.svcCtx.FavoriteMq, msg)
		return &favorite.FavoriteActionResp{}, nil
	}
	if in.ActionType == 2 {
		err := l.CancelFavoriteToRedis(in.UserId, in.VideoId)
		if err != nil {
			logx.Error(err)
			return nil, status.Error(500, err.Error())
		}

		msg, err := json.Marshal(in)
		if err != nil {
			logx.Error(err)
			return nil, status.Error(500, err.Error())
		}
		favoritemq.FavoriteActionSend(l.svcCtx.FavoriteMq, msg)
		return &favorite.FavoriteActionResp{}, nil
	}
	return nil, status.Error(500, errors.New("参数不合法").Error())
}

func (l *FavoriteActionLogic) SetFavoriteToRedis(uid int64, vid int64) error {
	uidStr := strconv.Itoa(int(uid))
	vidStr := strconv.Itoa(int(vid))
	key := prefix + vidStr
	isMember, err := l.svcCtx.Redis.SismemberCtx(l.ctx, key, uidStr)
	if err != nil {
		return err
	}
	if isMember {
		//logx.Error(fmt.Sprintf("user: %s favourite: %s conflict", uidStr, vidStr))
		return nil
	}
	_, err = l.svcCtx.Redis.SaddCtx(l.ctx, key, uidStr)
	if err != nil {
		logx.Error(fmt.Sprintf("user: %s set favourite to video: %s fail", uidStr, vidStr))
		return err
	}
	// 对于点赞的redis不设置过期时间，可以使用定时任务，定期更新redis的数据
	return nil
}

func (l *FavoriteActionLogic) CancelFavoriteToRedis(uid int64, vid int64) error {
	uidStr := strconv.Itoa(int(uid))
	vidStr := strconv.Itoa(int(vid))
	key := prefix + vidStr
	isMember, err := l.svcCtx.Redis.SismemberCtx(l.ctx, key, uidStr)
	if err != nil {
		return err
	}
	if !isMember {
		//logx.Error(fmt.Sprintf("user: %s favourite: %s not exist", uidStr, vidStr))
		//return errors.New("user: " + uidStr + " favourite : " + vidStr + " not exist")
		return nil
	}
	_, err = l.svcCtx.Redis.SremCtx(l.ctx, key, uidStr)
	if err != nil {
		logx.Error(fmt.Sprintf("user: %s cancel favourite to video: %s fail", uidStr, vidStr))
		return err
	}
	return nil
}
