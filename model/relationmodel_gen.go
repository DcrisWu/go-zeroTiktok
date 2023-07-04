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
	relationFieldNames          = builder.RawFieldNames(&Relation{})
	relationRows                = strings.Join(relationFieldNames, ",")
	relationRowsExpectAutoSet   = strings.Join(stringx.Remove(relationFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	relationRowsWithPlaceHolder = strings.Join(stringx.Remove(relationFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	relationModel interface {
		Insert(ctx context.Context, data *Relation) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Relation, error)
		FindOneByUserIdToUserId(ctx context.Context, userId int64, toUserId int64) (*Relation, error)
		Update(ctx context.Context, data *Relation) error
		Delete(ctx context.Context, id int64) error
	}

	defaultRelationModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Relation struct {
		Id        int64        `db:"id"`
		UserId    int64        `db:"user_id"`
		ToUserId  int64        `db:"to_user_id"`
		CreatedAt time.Time    `db:"created_at"`
		UpdatedAt time.Time    `db:"updated_at"`
		DeletedAt sql.NullTime `db:"deleted_at"`
	}
)

func newRelationModel(conn sqlx.SqlConn) *defaultRelationModel {
	return &defaultRelationModel{
		conn:  conn,
		table: "`relation`",
	}
}

func (m *defaultRelationModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultRelationModel) FindOne(ctx context.Context, id int64) (*Relation, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", relationRows, m.table)
	var resp Relation
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

func (m *defaultRelationModel) FindOneByUserIdToUserId(ctx context.Context, userId int64, toUserId int64) (*Relation, error) {
	var resp Relation
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `to_user_id` = ? limit 1", relationRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, toUserId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultRelationModel) Insert(ctx context.Context, data *Relation) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, relationRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.UserId, data.ToUserId, data.DeletedAt)
	return ret, err
}

func (m *defaultRelationModel) Update(ctx context.Context, newData *Relation) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, relationRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.UserId, newData.ToUserId, newData.DeletedAt, newData.Id)
	return err
}

func (m *defaultRelationModel) tableName() string {
	return m.table
}
