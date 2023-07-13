package pack

import (
	"context"
	"github.com/pkg/errors"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/relation-service/pb/relation"
	"gorm.io/gorm"
)

func RelationUser(ctx context.Context, DB *gorm.DB, u *db.User, fromId int64) (*relation.User, error) {
	if u == nil {
		return &relation.User{
			Name: "用户已注销",
		}, nil
	}

	followCount := int64(u.FollowingCount)
	followerCount := int64(u.FollowerCount)

	// true -> fromId 已关注了 u.Id, false-> fromId 未关注 u.Id
	isFollow := false
	rel, err := db.GetRelation(ctx, DB, fromId, int64(u.ID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if rel != nil {
		isFollow = true
	}
	return &relation.User{
		Id:            int64(u.ID),
		Name:          u.UserName,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
	}, nil
}

func RelationUsers(ctx context.Context, DB *gorm.DB, us []*db.User, fromId int64) ([]*relation.User, error) {
	users := make([]*relation.User, 0)
	for _, u := range us {
		user2, err := RelationUser(ctx, DB, u, fromId)
		if err != nil {
			return nil, err
		}
		if user2 != nil {
			users = append(users, user2)
		}
	}
	return users, nil
}

func FollowingList(ctx context.Context, DB *gorm.DB, vs []*db.Relation, fromId int64) ([]*relation.User, error) {
	users := make([]*db.User, 0)
	for _, v := range vs {
		user2, err := db.GetUserById(ctx, DB, int64(v.ToUserID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		users = append(users, user2)
	}
	return RelationUsers(ctx, DB, users, fromId)
}

func FollowerList(ctx context.Context, DB *gorm.DB, vs []*db.Relation, fromId int64) ([]*relation.User, error) {
	users := make([]*db.User, 0)
	for _, v := range vs {
		user2, err := db.GetUserById(ctx, DB, int64(v.UserID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		users = append(users, user2)
	}
	return RelationUsers(ctx, DB, users, fromId)
}
