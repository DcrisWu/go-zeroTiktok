package db

import (
	"context"
	"google.golang.org/grpc/status"
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

func CreateFavorite(ctx context.Context, DB *gorm.DB, uid int64, vid int64) error {
	// 事务
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 新增点赞数据
		user := new(User)
		if err := tx.WithContext(ctx).First(user, uid).Error; err != nil {
			return err
		}

		video := new(Video)
		if err := tx.WithContext(ctx).First(video, vid).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Model(&user).Association("FavoriteVideos").Append(video); err != nil {
			return err
		}

		// 2. 改变 video 表中的 favorite count
		res := tx.Model(video).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return status.Errorf(10001, "ErrDatabase")
		}
		return nil
	})
	return err
}

func CancelFavorite(ctx context.Context, DB *gorm.DB, uid int64, vid int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1.删除点赞数据
		user := new(User)
		if err := tx.WithContext(ctx).First(user, uid).Error; err != nil {
			return err
		}
		video, err := GetFavoriteRelation(ctx, DB, uid, vid)
		if err != nil {
			return err
		}

		err = tx.Unscoped().WithContext(ctx).Model(&user).Association("FavoriteVideos").Delete(video)
		if err != nil {
			return err
		}

		// 2.改变 video 表中的 favorite count
		res := tx.Model(video).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return status.Errorf(10001, "ErrDatabase")
		}
		return nil
	})
	return err
}

func FavoriteList(ctx context.Context, DB *gorm.DB, uid int64) ([]*Video, error) {
	user := new(User)
	if err := DB.WithContext(ctx).First(user, uid).Error; err != nil {
		return nil, err
	}

	var videos []*Video
	if err := DB.WithContext(ctx).Model(&user).Association("FavoriteVideos").Find(&videos); err != nil {
		return nil, err
	}
	return videos, nil
}
