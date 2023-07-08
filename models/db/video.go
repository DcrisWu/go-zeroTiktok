package db

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model
	UpdatedAt     time.Time `gorm:"column:update_time;not null;index:idx_update" `
	AuthorID      int       `gorm:"index:idx_authorid;not null"`
	PlayUrl       string    `gorm:"type:varchar(255);not null"`
	CoverUrl      string    `gorm:"type:varchar(255)"`
	FavoriteCount int       `gorm:"default:0"`
	CommentCount  int       `gorm:"default:0"`
	Title         string    `gorm:"type:varchar(50);not null"`
}

func (Video) TableName() string {
	return "video"
}

func GetVideos(ctx context.Context, DB *gorm.DB, limit int, latestTime int64) ([]*Video, error) {
	videos := make([]*Video, 0)

	if latestTime == 0 {
		latestTime = time.Now().UnixMilli()
	}
	conn := DB.WithContext(ctx)

	if err := conn.Limit(limit).Order("update_time desc").Find(&videos, "update_time < ?", time.UnixMilli(latestTime)).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func CreateVideo(ctx context.Context, DB *gorm.DB, video *Video) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(video).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// PublishList 获取 authorId 发布过的所有视频
func PublishList(ctx context.Context, DB *gorm.DB, authorId int64) ([]*Video, error) {
	var pubList []*Video
	if err := DB.WithContext(ctx).Model(&Video{}).Where(&Video{AuthorID: int(authorId)}).Find(&pubList).Error; err != nil {
		return nil, err
	}
	return pubList, nil
}
