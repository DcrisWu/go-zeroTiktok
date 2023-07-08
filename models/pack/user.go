package pack

import (
	"context"
	"github.com/pkg/errors"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/user-service/pb/user"
	"gorm.io/gorm"
)

// User 将 db.User 结构体转为 user.UserInfo
func User(ctx context.Context, DB *gorm.DB, u *db.User, fromId int64) (*user.UserInfo, error) {
	if u == nil {
		return &user.UserInfo{
			Name: "用户已注销",
		}, nil
	}

	followCount := int64(u.FollowingCount)
	followerCount := int64(u.FollowerCount)

	// true -> fromId 已关注了 u.Id, false-> fromId 未关注 u.Id
	isFollow := false
	relation, err := db.GetRelation(ctx, DB, fromId, int64(u.ID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if relation != nil {
		isFollow = true
	}
	return &user.UserInfo{
		Id:            int64(u.ID),
		Name:          u.UserName,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
	}, nil
}
