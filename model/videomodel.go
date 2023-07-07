package model

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zeroTiktok/model/db"
	"time"
)

var _ VideoModel = (*customVideoModel)(nil)

type (
	// VideoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideoModel.
	VideoModel interface {
		videoModel
		WithTx(session sqlx.Session) VideoModel
		GetVideosByAuthorId(ctx context.Context, AuthorId, limit, offset int64) ([]*Video, error)
		GetVideosOrderByTime(ctx context.Context, limit, offset int64, LastTime time.Time) ([]*Video, error)
		IncrCommentCount(ctx context.Context, num, videoId int64) error
	}

	customVideoModel struct {
		*defaultVideoModel
		genericModel db.GenericModelConnI[Video]
	}
)

// NewVideoModel returns a model for the database table.
func NewVideoModel(conn sqlx.SqlConn) VideoModel {
	return &customVideoModel{
		defaultVideoModel: newVideoModel(conn),
		genericModel:      db.NewGenericModelConn[Video](conn, "video"),
	}
}

func (c *customVideoModel) WithTx(session sqlx.Session) VideoModel {
	return &customVideoModel{
		defaultVideoModel: c.defaultVideoModel,
		genericModel:      c.genericModel.WithTx(session),
	}
}

func (c *customVideoModel) GetVideosByAuthorId(ctx context.Context, AuthorId, limit, offset int64) ([]*Video, error) {
	where := []goqu.Expression{}
	where = append(where, goqu.I("author_id").Eq(AuthorId))
	return c.genericModel.List(ctx, where, []exp.OrderedExpression{goqu.I("created_at").Desc()}, nil, limit, offset)
}

func (c *customVideoModel) GetVideosOrderByTime(ctx context.Context, limit, offset int64, LastTime time.Time) ([]*Video, error) {
	where := []goqu.Expression{}
	where = append(where, goqu.I("created_at").Lt(LastTime))
	return c.genericModel.List(ctx, where, []exp.OrderedExpression{goqu.I("created_at").Desc()}, nil, limit, offset)
}

func (c *customVideoModel) IncrCommentCount(ctx context.Context, num, videoId int64) error {
	newValue := goqu.Record{
		"comment_count": goqu.L(fmt.Sprintf("%s+(%d)", "comment_count", num)),
	}
	where := []goqu.Expression{goqu.C("id").Eq(videoId)}
	_, err := c.genericModel.Update(ctx, newValue, where)
	return err
}
