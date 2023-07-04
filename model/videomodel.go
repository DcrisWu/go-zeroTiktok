package model

import (
	"context"
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
		GetVideosByAuthorId(ctx context.Context, AuthorId, limit, offset int64) ([]*Video, error)
		GetVideosOrderByTime(ctx context.Context, limit, offset int64, LastTime time.Time) ([]*Video, error)
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
