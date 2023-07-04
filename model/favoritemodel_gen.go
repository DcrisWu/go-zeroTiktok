// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	favoriteFieldNames          = builder.RawFieldNames(&Favorite{})
	favoriteRows                = strings.Join(favoriteFieldNames, ",")
	favoriteRowsExpectAutoSet   = strings.Join(stringx.Remove(favoriteFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	favoriteRowsWithPlaceHolder = strings.Join(stringx.Remove(favoriteFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	favoriteModel interface {
		Insert(ctx context.Context, data *Favorite) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Favorite, error)
		FindOneByUserIdVideoId(ctx context.Context, userId int64, videoId int64) (*Favorite, error)
		Update(ctx context.Context, data *Favorite) error
		Delete(ctx context.Context, id int64) error
	}

	defaultFavoriteModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Favorite struct {
		Id        int64     `db:"id"`
		UserId    int64     `db:"user_id"`
		VideoId   int64     `db:"video_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
)

func newFavoriteModel(conn sqlx.SqlConn) *defaultFavoriteModel {
	return &defaultFavoriteModel{
		conn:  conn,
		table: "`favorite`",
	}
}

func (m *defaultFavoriteModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultFavoriteModel) FindOne(ctx context.Context, id int64) (*Favorite, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", favoriteRows, m.table)
	var resp Favorite
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFavoriteModel) FindOneByUserIdVideoId(ctx context.Context, userId int64, videoId int64) (*Favorite, error) {
	var resp Favorite
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `video_id` = ? limit 1", favoriteRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, videoId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFavoriteModel) Insert(ctx context.Context, data *Favorite) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, favoriteRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.VideoId)
	return ret, err
}

func (m *defaultFavoriteModel) Update(ctx context.Context, newData *Favorite) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, favoriteRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.UserId, newData.VideoId, newData.Id)
	return err
}

func (m *defaultFavoriteModel) tableName() string {
	return m.table
}
