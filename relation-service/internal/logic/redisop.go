package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/utils"
	"gorm.io/gorm"
	"strconv"
)

func addRedisFollowList(ctx context.Context, DB *gorm.DB, redis *redis.Redis, userId int64, toUserId int64) error {
	uidStr := strconv.Itoa(int(userId))
	toUserIdStr := strconv.Itoa(int(toUserId))
	loadFollowingListToRedis(ctx, DB, redis, userId)
	key := "following:" + uidStr
	_, err := redis.SaddCtx(ctx, key, toUserIdStr)
	if err != nil {
		return err
	}
	err = redis.ExpireCtx(ctx, key, utils.RedisExpireTime)
	if err != nil {
		return err
	}
	return nil
}

func addRedisFollowerList(ctx context.Context, DB *gorm.DB, redis *redis.Redis, userId int64, toUserId int64) error {
	uidStr := strconv.Itoa(int(userId))
	toUserIdStr := strconv.Itoa(int(toUserId))
	loadFollowersListToRedis(ctx, DB, redis, toUserId)
	key := "follower:" + toUserIdStr
	_, err := redis.Sadd(key, uidStr)
	if err != nil {
		return err
	}
	err = redis.ExpireCtx(ctx, key, utils.RedisExpireTime)
	if err != nil {
		return err
	}
	return nil
}

func rmRedisFollowList(ctx context.Context, DB *gorm.DB, redis *redis.Redis, userId int64, toUserId int64) error {
	uidStr := strconv.Itoa(int(userId))
	toUserIdStr := strconv.Itoa(int(toUserId))
	loadFollowingListToRedis(ctx, DB, redis, userId)
	key := "following:" + uidStr
	// redis中不存在该用户的关注列表，不对redis进行操作
	_, err := redis.SremCtx(ctx, key, toUserIdStr)
	if err != nil {
		return err
	}
	err = redis.ExpireCtx(ctx, key, utils.RedisExpireTime)
	if err != nil {
		return err
	}
	return nil
}

func rmRedisFollowerList(ctx context.Context, DB *gorm.DB, redis *redis.Redis, userId int64, toUserId int64) error {
	uidStr := strconv.Itoa(int(userId))
	toUserIdStr := strconv.Itoa(int(toUserId))
	loadFollowersListToRedis(ctx, DB, redis, toUserId)
	key := "follower:" + toUserIdStr
	// redis中不存在该用户的关注列表，不对redis进行操作
	_, err := redis.SremCtx(ctx, key, uidStr)
	if err != nil {
		return err
	}
	err = redis.ExpireCtx(ctx, key, utils.RedisExpireTime)
	if err != nil {
		return err
	}
	return nil
}

func loadFollowingListToRedis(ctx context.Context, DB *gorm.DB, redis *redis.Redis, userId int64) error {
	uidStr := strconv.Itoa(int(userId))
	key := "following:" + uidStr
	exist, err := redis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Error("加载关注列表到redis失败")
		return err
	}
	if exist {
		redis.ExpireCtx(ctx, key, utils.RedisExpireTime)
		return nil
	}
	// 从数据库中加载关注列表到redis
	ids, err := db.GetFollowingIdByUserId(ctx, DB, int(userId))
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	// 将关注列表加载到redis
	for _, id := range ids {
		_, err := redis.SaddCtx(ctx, key, strconv.Itoa(int(id)))
		if err != nil {
			return err
		}
	}
	err = redis.ExpireCtx(ctx, key, utils.RedisExpireTime)
	if err != nil {
		return err
	}
	return nil
}

func loadFollowersListToRedis(ctx context.Context, DB *gorm.DB, redis *redis.Redis, toUserid int64) error {
	uidStr := strconv.Itoa(int(toUserid))
	key := "follower:" + uidStr
	exist, err := redis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Error("加载粉丝列表到redis失败")
		return err
	}
	if exist {
		redis.ExpireCtx(ctx, key, utils.RedisExpireTime)
		return nil
	}
	// 从数据库中加载粉丝列表到redis
	ids, err := db.GetFollowerIdByUserId(ctx, DB, int(toUserid))
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	// 将粉丝列表加载到redis
	for _, id := range ids {
		_, err := redis.SaddCtx(ctx, key, strconv.Itoa(int(id)))
		if err != nil {
			return err
		}
	}
	err = redis.ExpireCtx(ctx, key, utils.RedisExpireTime)
	if err != nil {
		return err
	}
	return nil
}
