package db

import (
	"context"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName       string  `gorm:"index:idx_username,unique;type:varchar(40);not null" json:"user_name"`
	Password       string  `gorm:"type:varchar(256);not null" json:"password"`
	FavoriteVideos []Video `gorm:"many2many:user_favorite_videos" json:"favorite_videos"`
	FollowingCount int     `gorm:"default:0" json:"following_count"`
	FollowerCount  int     `gorm:"default:0" json:"follower_count"`
}

func (User) TableName() string {
	return "user"
}

func QueryUser(ctx context.Context, DB *gorm.DB, userName string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(ctx).Where("user_name = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func CreateUser(ctx context.Context, DB *gorm.DB, users []*User) error {
	return DB.WithContext(ctx).Create(users).Error
}

func GetUserById(ctx context.Context, DB *gorm.DB, userId int64) (*User, error) {
	res := new(User)
	if err := DB.WithContext(ctx).First(&res, userId).Error; err != nil {
		return nil, err
	}
	return res, nil
}
