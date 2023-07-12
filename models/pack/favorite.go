package pack

import (
	"context"
	"errors"
	"go-zeroTiktok/favorite-service/pb/favorite"
	"go-zeroTiktok/models/db"
	"gorm.io/gorm"
)

// FeedVideo pack feed info
func FavoriteVideo(ctx context.Context, DB *gorm.DB, v *db.Video, fromID int64) (*favorite.Video, error) {
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

	author := &favorite.User{
		Id:            u.Id,
		Name:          u.Name,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		IsFollow:      u.IsFollow,
	}

	return &favorite.Video{
		Id:            int64(v.ID),
		Author:        author,
		PlayUrl:       v.PlayUrl,
		CoverUrl:      v.CoverUrl,
		FavoriteCount: favoriteCount,
		CommentCount:  commentCount,
		Title:         v.Title,
	}, nil
}

func FavoriteVideos(ctx context.Context, DB *gorm.DB, vs []*db.Video, fromID int64) ([]*favorite.Video, error) {
	videos := make([]*favorite.Video, 0)
	for _, v := range vs {
		video2, err := FavoriteVideo(ctx, DB, v, fromID)
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
			video2.IsFavorite = flag
			videos = append(videos, video2)
		}
	}
	return videos, nil
}
