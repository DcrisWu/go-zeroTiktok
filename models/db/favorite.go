package db

import (
	"context"
	"gorm.io/gorm"
)

// GetFavoriteRelation 获取 favorite video info
func GetFavoriteRelation(ctx context.Context, DB *gorm.DB, uid int64, vid int64) (*Video, error) {
	user := new(User)
	if err := DB.WithContext(ctx).First(user, uid).Error; err != nil {
		return nil, err
	}

	video := new(Video)
	if err := DB.WithContext(ctx).Model(&user).Association("FavoriteVideos").Find(&video, vid); err != nil {
		return nil, err
	}
	return video, nil
}
