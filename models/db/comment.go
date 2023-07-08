package db

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserId  int    `gorm:"index:idx_userid;not null"`
	VideoId int    `gorm:"index:idx_videoid;not null"`
	Content string `gorm:"type:varchar(255);not null"`
}

func (Comment) TableName() string {
	return "comment"
}

func NewComment(ctx context.Context, DB *gorm.DB, comment *Comment) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 事务执行
		// 1. 新增评论数据
		err := tx.Create(comment).Error
		if err != nil {
			return err
		}
		// 2.video表格中对应的comment_count增加
		res := tx.Model(&Video{}).Where("ID = ?", comment.VideoId).Update("comment_count", gorm.Expr("comment_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.Errorf("数据库错误")
		}
		return nil
	})
	return err
}

func DeleteComment(ctx context.Context, DB *gorm.DB, commentId int64, vid int64) error {
	// 事务
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		comment := new(Comment)
		if err := tx.First(&comment, commentId).Error; err != nil {
			return err
		}

		// 1.删除评论数据
		if err := tx.Unscoped().Delete(&comment).Error; err != nil {
			return err
		}

		//2.改变video表格中的video count
		res := tx.Model(&Video{}).Where("ID = ?", vid).Update("comment_count", gorm.Expr("comment_count - ?", 1))
		if res.RowsAffected != 1 {
			return errors.Errorf("数据库发生错误")
		}
		return nil
	})
	return err
}

func GetVideoComments(ctx context.Context, DB *gorm.DB, vid int64) ([]*Comment, error) {
	var comments []*Comment
	if err := DB.WithContext(ctx).Model(&Comment{}).Where(&Comment{VideoId: int(vid)}).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
