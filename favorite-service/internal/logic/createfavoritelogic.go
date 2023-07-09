package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zeroTiktok/models/db"
	"google.golang.org/grpc/status"

	"go-zeroTiktok/favorite-service/internal/svc"
	"go-zeroTiktok/favorite-service/pb/favorite"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFavoriteLogic {
	return &CreateFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateFavoriteLogic) CreateFavorite(in *favorite.CreateFavoriteReq) (*favorite.CreateFavoriteResp, error) {
	key := fmt.Sprintf("write:favorite:video:%d:%d", in.UserId, in.VideoId)
	err := l.svcCtx.Redis.Expireat(key, int64(3600))
	lock := redis.NewRedisLock(l.svcCtx.Redis, key)
	lock.SetExpire(3600)
	success, err := lock.Acquire()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	if !success {
		return nil, status.Error(400, "已经点赞过了")
	}
	go func(prefix string, userId int64, videoId int64, redisLock *redis.RedisLock) {
		// todo
		// 错误处理：休眠一段时间后重试，重试次数限制 5 次，超过 5 次则放弃，同时删除redis key
		err = db.CreateFavorite(l.ctx, l.svcCtx.DB, userId, videoId)
		if err != nil {
			logx.Error(err)
			return
		}
		// 释放锁，释放失败也没关系，锁会在过期后自动释放
		redisLock.Release()
		favCount := fmt.Sprintf("like:video:%d", videoId)
		// 为了防止用户重复点赞，将用户点赞的视频id写入redis，设置过期时间为 1 小时
		userVideo := fmt.Sprintf("like:video-user:%d:%d", videoId, userId)
		l.svcCtx.Redis.SetnxEx(userVideo, "0", 3600)
		// 下单成功后，将视频点赞数写进redis，写入失败也没关系，可以使用定时任务将数据库中的数据写入redis
		l.svcCtx.Redis.Incrby(favCount, 1)
	}("favorite:video", in.UserId, in.VideoId, lock)
	// 在写入数据库操作成功前，先返回成功，避免用户重复点击，再后台刷新数据库
	// 返回成功后，前端可以暂时显示点赞成功，并设置超时时间，超时后请求接口刷新数据
	// 后台刷新成功，前端可以获取到正确的数据，刷新失败，前端就可以显示点赞失败，用户可以重新点击
	return &favorite.CreateFavoriteResp{
		IsSuccess: true,
	}, nil
}

//func flushDataBase()  {
//
//}
