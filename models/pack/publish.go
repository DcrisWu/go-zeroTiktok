package pack

import (
	"context"
	"errors"
	"go-zeroTiktok/models/db"
	"go-zeroTiktok/publish-service/pb/publish"
	"gorm.io/gorm"
)

// PublishVideo pack feed info
func PublishVideo(ctx context.Context, DB *gorm.DB, v *db.Video, fromID int64) (*publish.Video, error) {
	if v == nil {
		return nil, nil
	}
	user, err := db.GetUserById(ctx, DB, int64(v.AuthorID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	u, err := User(ctx, DB, user, fromID)
	if err != nil {
		return nil, err
	}
	favoriteCount := int64(v.FavoriteCount)
	commentCount := int64(v.CommentCount)

	author := &publish.User{
		Id:            u.Id,
		Name:          u.Name,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		IsFollow:      u.IsFollow,
	}

	return &publish.Video{
		Id:           int64(v.ID),
		Author:       author,
		PlayUrl:      v.PlayUrl,
		CoverUrl:     v.CoverUrl,
		FavorCount:   favoriteCount,
		CommentCount: commentCount,
		Title:        v.Title,
	}, nil
}

// PublishVideos pack list of video info
func PublishVideos(ctx context.Context, DB *gorm.DB, vs []*db.Video, fromID int64) ([]*publish.Video, error) {
	videos := make([]*publish.Video, 0)
	for _, v := range vs {
		video2, err := PublishVideo(ctx, DB, v, fromID)
		if err != nil {
			return nil, err
		}

		if video2 != nil {
			flag := false
			if fromID != 0 {
				results, err := db.GetFavoriteRelation(ctx, DB, fromID, int64(v.ID))
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				} else if errors.Is(err, gorm.ErrRecordNotFound) {
					flag = false
				} else if results != nil && results.AuthorID != 0 {
					flag = true
				}
			}
			video2.IsFavor = flag
			videos = append(videos, video2)
		}
	}
	return videos, nil
}
