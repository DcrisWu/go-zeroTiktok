package model

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zeroTiktok/model/db"
)

var _ CommentModel = (*customCommentModel)(nil)

type (
	// CommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentModel.
	CommentModel interface {
		commentModel
		WithTx(session sqlx.Session) CommentModel
		DeleteByCommentId(ctx context.Context, commentId int64) error
		ListCommentByVideoId(ctx context.Context, videoId int64) ([]*Comment, error)
	}

	customCommentModel struct {
		*defaultCommentModel
		genericModel db.GenericModelConnI[Video]
	}
)

// NewCommentModel returns a model for the database table.
func NewCommentModel(conn sqlx.SqlConn) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn),
		genericModel:        db.NewGenericModelConn[Video](conn, "comment"),
	}
}

func (c *customCommentModel) WithTx(session sqlx.Session) CommentModel {
	return &customCommentModel{
		defaultCommentModel: c.defaultCommentModel,
		genericModel:        c.genericModel.WithTx(session),
	}
}

func (c *customCommentModel) DeleteByCommentId(ctx context.Context, commentId int64) error {
	var where []goqu.Expression
	where = append(where, goqu.I("video_id").Eq(commentId))
	_, err := c.genericModel.Delete(ctx, where)
	return err
}

func (c *customCommentModel) ListCommentByVideoId(ctx context.Context, videoId int64) ([]*Comment, error) {
	query := fmt.Sprintf("select %s from %s where `video_id` = ? order by 'created_at'", commentRows, c.table)
	var resp []*Comment
	err := c.conn.QueryRows(ctx, query, videoId, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
