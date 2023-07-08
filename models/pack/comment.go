package pack

import (
	"context"
	"errors"
	"go-zeroTiktok/comment-service/pb/comment"
	"go-zeroTiktok/models/db"
	"gorm.io/gorm"
)

func Comments(ctx context.Context, DB *gorm.DB, vs []*db.Comment, fromID int64) ([]*comment.CommentInfo, error) {
	comments := make([]*comment.CommentInfo, 0)
	for _, v := range vs {
		u, err := db.GetUserById(ctx, DB, int64(v.UserId))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		packUser, err := CommentUser(ctx, DB, u, fromID)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment.CommentInfo{
			CommentId: int64(v.ID),
			User:      packUser,
			Content:   v.Content,
			CreatedAt: v.CreatedAt.Format("2023-01-02"),
		})
	}
	return comments, nil
}

func CommentUser(ctx context.Context, DB *gorm.DB, u *db.User, fromId int64) (*comment.User, error) {
	if u == nil {
		return &comment.User{
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
	return &comment.User{
		Id:            int64(u.ID),
		Name:          u.UserName,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
	}, nil
}
