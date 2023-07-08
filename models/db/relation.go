package db

import (
	"context"
	"gorm.io/gorm"
)

// Relation Gorm data structure
// Relation 既属于 关注者 也属于 被关注者
type Relation struct {
	gorm.Model
	UserID   int `gorm:"index:idx_userid,unique;not null"`
	ToUserID int `gorm:"index:idx_userid,unique;index:idx_userid_to;not null"`
}

func (Relation) TableName() string {
	return "relation"
}

func GetRelation(ctx context.Context, DB *gorm.DB, uid int64, tid int64) (*Relation, error) {
	relation := new(Relation)

	if err := DB.WithContext(ctx).First(&relation, "user_id = ? and to_user_id = ?", uid, tid).Error; err != nil {
		return nil, err
	}
	return relation, nil
}
