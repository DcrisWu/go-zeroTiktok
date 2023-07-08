package model

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
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
		InsertComment(ctx context.Context, data *Comment) error
		DeleteByCommentId(ctx context.Context, commentId int64) error
		ListCommentByVideoId(ctx context.Context, videoId int64) ([]*Comment, error)
	}

	customCommentModel struct {
		*defaultCommentModel
		genericModel db.GenericModelConnI[Comment]
	}
)

// NewCommentModel returns a model for the database table.
func NewCommentModel(conn sqlx.SqlConn) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn),
		genericModel:        db.NewGenericModelConn[Comment](conn, "comment"),
	}
}

func (c *customCommentModel) WithTx(session sqlx.Session) CommentModel {
	return &customCommentModel{
		defaultCommentModel: c.defaultCommentModel,
		genericModel:        c.genericModel.WithTx(session),
	}
}

func (c *customCommentModel) InsertComment(ctx context.Context, data *Comment) error {
	_, err := c.genericModel.Insert(ctx, data)
	return err
}

func (c *customCommentModel) DeleteByCommentId(ctx context.Context, commentId int64) error {
	var where []goqu.Expression
	where = append(where, goqu.I("video_id").Eq(commentId))
	_, err := c.genericModel.Delete(ctx, where)
	return err
}

func (c *customCommentModel) ListCommentByVideoId(ctx context.Context, videoId int64) ([]*Comment, error) {
	where := goqu.And(
		goqu.I("video_id").Eq(videoId),
	)
	var orderBy []exp.OrderedExpression
	orderBy = append(orderBy, goqu.I("id").Desc())
	return c.genericModel.List(ctx, where.Expressions(), orderBy, nil, 0, 0)
}
