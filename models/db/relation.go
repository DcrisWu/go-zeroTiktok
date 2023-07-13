package db

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Relation Gorm data structure
// UserId 关注 ToUserId
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

//func GetFollowingIdByUserId(ctx context.Context, DB *gorm.DB, userId int) ([]int64, error) {
//	var following []int64
//	err := DB.WithContext(ctx).Select("to_user_id").Where("user_id = ?", userId).Find(&following).Error
//	if err != nil {
//		return nil, err
//	}
//	return following, err
//}
//
//func GetFollowerIdByUserId(ctx context.Context, DB *gorm.DB, userId int) ([]int64, error) {
//	var follower []int64
//	err := DB.WithContext(ctx).Select("user_id").Where("to_user_id = ?", userId).Find(&follower).Error
//	if err != nil {
//		return nil, err
//	}
//	return follower, err
//}

func CreateRelation(ctx context.Context, DB *gorm.DB, uid int64, tid int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 新增关注数据
		err := tx.Create(&Relation{UserID: int(uid), ToUserID: int(tid)}).Error
		if err != nil {
			return err
		}

		// 2.改变 user 表中的 following count
		res := tx.Model(new(User)).Where("ID = ?", uid).Update("following_count", gorm.Expr("following_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("ErrDatabase")
		}

		// 3.改变 user 表中的 follower count
		res = tx.Model(new(User)).Where("ID = ?", tid).Update("follower_count", gorm.Expr("follower_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("ErrDatabase")
		}

		return nil
	})
	return err
}

func CancelRelation(ctx context.Context, DB *gorm.DB, uid int64, tid int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		relation := new(Relation)
		if err := tx.Where("user_id = ? AND to_user_id=?", uid, tid).First(&relation).Error; err != nil {
			return err
		}
		// 1. 删除关注数据
		err := tx.Unscoped().Delete(&relation).Error
		if err != nil {
			return err
		}
		// 2.改变 user 表中的 following count
		res := tx.Model(new(User)).Where("ID = ?", uid).Update("following_count", gorm.Expr("following_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("ErrDatabase")
		}

		// 3.改变 user 表中的 follower count
		res = tx.Model(new(User)).Where("ID = ?", tid).Update("follower_count", gorm.Expr("follower_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("ErrDatabase")
		}

		return nil
	})
	return err
}

func ListFollowing(ctx context.Context, DB *gorm.DB, uid int64) ([]*Relation, error) {
	var relationList []*Relation
	err := DB.WithContext(ctx).Where("user_id = ?", uid).Find(&relationList).Error
	if err != nil {
		return nil, err
	}
	return relationList, nil
}

func ListFollower(ctx context.Context, DB *gorm.DB, tid int64) ([]*Relation, error) {
	var relationList []*Relation
	err := DB.WithContext(ctx).Where("to_user_id = ?", tid).Find(&relationList).Error
	if err != nil {
		return nil, err
	}
	return relationList, nil
}
