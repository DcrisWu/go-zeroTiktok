package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"reflect"
	"strings"

	goqumysql "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

func init() {
	mysqloptions := goqumysql.DialectOptionsV8()
	// INSERT IGNORE INTO `user_point` (`uid`, `point`, `test`) VALUES (13, 10, '1111111111111111') ON DUPLICATE KEY UPDATE
	// INSERT INTO `user_point` (`uid`, `point`, `test`) VALUES (1, 10, '1111111111111111') ON DUPLICATE KEY UPDATE
	// `SupportsInsertIgnoreSyntax` 打开，会多一个 IGNORE，会隐藏其他非主键冲突的异常 https://stackoverflow.com/questions/2366813/on-duplicate-key-ignore
	mysqloptions.SupportsInsertIgnoreSyntax = false
	goqu.RegisterDialect("bblmysql8", mysqloptions)
}

var ErrNotFound = sqlx.ErrNotFound

type IScroller interface {
	GetCursorOrderExp() exp.OrderedExpression
	GetCursorWhereExp() exp.Expression
	GetLimit() uint
}

type idScroller struct {
	limit     int64
	before    *int64
	after     *int64
	asc       bool
	cursorVal int64
	isUp      bool // 向上滚动还是向下滚动 true 表示是向上滚动
}

func NewIdScroller(limit int64, before, after *int64, asc bool) IScroller {
	var cursorVal int64
	var up bool
	if before != nil {
		cursorVal = *before
		up = true
	}
	if after != nil {
		cursorVal = *after
	}

	return &idScroller{
		limit:     limit,
		before:    before,
		after:     after,
		asc:       asc,
		cursorVal: cursorVal,
		isUp:      up,
	}
}

func (i *idScroller) GetCursorOrderExp() exp.OrderedExpression {
	cursorOrder := goqu.I("id").Desc()
	if i.asc {
		cursorOrder = goqu.I("id").Asc()
	}
	return cursorOrder
}

func (i *idScroller) GetCursorWhereExp() exp.Expression {
	/*
		limit=3

		Id = [1,2,3,4,5,6,7,8,9,10]
		Curr = [5,6,7]
		往上滚: before=5&asc=true -> [2,3,4]  id < before orderby id asc
		往下滚：after=7&asc=true -> [8,9,10]  id > after orderby id asc

		Id = [1,2,3,4,5,6,7,8,9,10]
		Curr = [7,6,5]
		往上滚: before=7&asc=false -> [10,9,8] id > before orderby id desc
		往下滚：after=5&asc=false -> [4,3,2] id < after orderby id  desc
	*/

	c := goqu.C("id")
	if i.cursorVal == 0 {
		return goqu.L("1=1")
	}
	if i.isUp {
		if i.asc {
			return c.Lt(i.cursorVal)
		}
		return c.Gt(i.cursorVal)
	} else {
		if i.asc {
			return c.Gt(i.cursorVal)
		}
		return c.Lt(i.cursorVal)
	}
}
func (i *idScroller) GetLimit() uint {
	return uint(i.limit)
}

type GenericModelConnI[T any] interface {
	One(ctx context.Context, where []goqu.Expression) (*T, error)
	ListAndCount(ctx context.Context, where []goqu.Expression, orderBy []exp.OrderedExpression, columns []any, limit, offset int64) ([]*T, int64, error)
	List(ctx context.Context, where []goqu.Expression, orderBy []exp.OrderedExpression, columns []any, limit, offset int64) ([]*T, error)
	Scroll(ctx context.Context, where []goqu.Expression, columns []any, scroller IScroller) ([]*T, error)
	Update(ctx context.Context, newValue goqu.Record, where []goqu.Expression) (sql.Result, error)
	Count(ctx context.Context, where []goqu.Expression) (int64, error)
	Sum(ctx context.Context, col exp.IdentifierExpression, where []goqu.Expression) (int64, error)
	InsertIgnoreConflict(ctx context.Context, model *T) (sql.Result, error)

	UpdateTx(ctx context.Context, sess sqlx.Session, newValue goqu.Record, where []goqu.Expression) (sql.Result, error)
	InsertTx(ctx context.Context, sess sqlx.Session, model *T, ignoreConfilict bool) (sql.Result, error)
	SelectOneForUpdateTx(ctx context.Context, sess sqlx.Session, where []goqu.Expression, block bool) (*T, error)
	SumTx(ctx context.Context, sess sqlx.Session, col exp.IdentifierExpression, where []goqu.Expression) (int64, error)
	ListTx(ctx context.Context, sess sqlx.Session, where []goqu.Expression, orderBy []exp.OrderedExpression, columns []any, limit, offset int64) ([]*T, error)
	CountTx(ctx context.Context, sess sqlx.Session, where []goqu.Expression) (int64, error)
	OneTx(ctx context.Context, sess sqlx.Session, where []goqu.Expression) (*T, error)

	GenUpdateSQL(newValue goqu.Record, where []goqu.Expression) (sql string, params []interface{}, err error)
	GenInsertSQL(model *T, onConflictDoNothing bool) (sql string, params []interface{}, err error)
	GenSelectForUpdateSQL(where []goqu.Expression, Block bool) (sql string, params []interface{}, err error)
}

type genericModelConn[T any] struct {
	conn      sqlx.SqlConn
	tableName string
	dialect   goqu.DialectWrapper
}

func NewGenericModelConn[T any](conn sqlx.SqlConn, tableName string) GenericModelConnI[T] {
	return &genericModelConn[T]{
		conn:      conn,
		tableName: tableName,
		dialect:   goqu.Dialect("bblmysql8"),
	}
}

func (g *genericModelConn[T]) list(ctx context.Context, conn sqlx.Session, where []goqu.Expression, orderBy []exp.OrderedExpression, columns []any, limit, offset int64) ([]*T, error) {
	partial := false
	selectDs := g.dialect.From(g.tableName).Where(where...).Limit(uint(limit)).Offset(uint(offset)).Order(orderBy...)
	if len(columns) != 0 {
		selectDs = selectDs.Select(columns...)
		partial = true
	}

	listQuery, listArgs, err := selectDs.Prepared(true).ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "select dataset prepared error")
	}

	entities := make([]*T, 0)
	if partial {
		err = conn.QueryRowsPartialCtx(ctx, &entities, listQuery, listArgs...)
	} else {
		err = conn.QueryRowsCtx(ctx, &entities, listQuery, listArgs...)
	}

	return entities, err
}

func (g *genericModelConn[T]) count(ctx context.Context, conn sqlx.Session, where []goqu.Expression) (int64, error) {
	countDs := g.dialect.From(g.tableName).Where(where...).Select(goqu.COUNT("*").As("total"))
	countQuery, countArgs, err := countDs.Prepared(true).ToSQL()
	if err != nil {
		return 0, errors.Wrapf(err, "count dataset prepared error")
	}

	var total int64
	err = conn.QueryRowCtx(ctx, &total, countQuery, countArgs...)
	return total, err
}

func (g *genericModelConn[T]) sum(ctx context.Context, conn sqlx.Session, col exp.IdentifierExpression, where []goqu.Expression) (int64, error) {
	sumDs := g.dialect.From(g.tableName).Where(where...).Select(goqu.COALESCE(goqu.SUM(col), 0).As("sum"))
	countQuery, countArgs, err := sumDs.Prepared(true).ToSQL()
	if err != nil {
		return 0, errors.Wrapf(err, "count dataset prepared error")
	}

	var s int64
	err = conn.QueryRowCtx(ctx, &s, countQuery, countArgs...)
	return s, err
}

func (g *genericModelConn[T]) one(ctx context.Context, conn sqlx.Session, where []goqu.Expression) (*T, error) {
	selectDs := g.dialect.From(g.tableName).Where(where...).Limit(1)
	query, args, err := selectDs.Prepared(true).ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "one prepare error")
	}

	var resp T
	err = conn.QueryRowCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (g *genericModelConn[T]) insert(ctx context.Context, conn sqlx.Session, model *T, ignoreConflict bool) (sql.Result, error) {
	sql, args, err := g.GenInsertSQL(model, ignoreConflict)
	if err != nil {
		return nil, errors.Wrapf(err, "gen insert sql error")
	}

	return conn.ExecCtx(ctx, sql, args...)

}

func (g *genericModelConn[T]) InsertIgnoreConflict(ctx context.Context, model *T) (sql.Result, error) {
	return g.insert(ctx, g.conn, model, true)
}

func (g *genericModelConn[T]) InsertTx(ctx context.Context, sess sqlx.Session, model *T, ignoreConflict bool) (sql.Result, error) {
	return g.insert(ctx, sess, model, ignoreConflict)
}

func (g *genericModelConn[T]) One(ctx context.Context, where []goqu.Expression) (*T, error) {
	return g.one(ctx, g.conn, where)
}

func (g *genericModelConn[T]) OneTx(ctx context.Context, sess sqlx.Session, where []goqu.Expression) (*T, error) {
	return g.one(ctx, sess, where)
}

func (g *genericModelConn[T]) ListAndCount(ctx context.Context, where []goqu.Expression, orderBy []exp.OrderedExpression, columns []any, limit, offset int64) (ret []*T, total int64, err error) {
	ret, err = g.List(ctx, where, orderBy, columns, limit, offset)
	if err != nil {
		return
	}

	total, err = g.Count(ctx, where)
	return
}

func (g *genericModelConn[T]) Scroll(ctx context.Context, where []goqu.Expression, columns []any, scroller IScroller) (ret []*T, err error) {
	selectWhere := goqu.And(append(where, scroller.GetCursorWhereExp())...)
	selectDs := g.dialect.From(g.tableName).Where(selectWhere.Expressions()...).Order(scroller.GetCursorOrderExp()).Limit(scroller.GetLimit())
	partial := false
	if len(columns) != 0 {
		selectDs = selectDs.Select(columns...)
		partial = true
	}

	listQuery, listArgs, err := selectDs.Prepared(true).ToSQL()
	if err != nil {
		err = errors.Wrapf(err, "scroll select prepare error")
		return
	}

	entities := make([]*T, 0)
	if partial {
		err = g.conn.QueryRowsPartialCtx(ctx, &entities, listQuery, listArgs...)
	} else {
		err = g.conn.QueryRowsCtx(ctx, &entities, listQuery, listArgs...)
	}

	return entities, err
}

func (g *genericModelConn[T]) Count(ctx context.Context, where []goqu.Expression) (total int64, err error) {
	return g.count(ctx, g.conn, where)
}

func (g *genericModelConn[T]) CountTx(ctx context.Context, sess sqlx.Session, where []goqu.Expression) (total int64, err error) {
	return g.count(ctx, sess, where)
}

func (g *genericModelConn[T]) List(ctx context.Context, where []goqu.Expression, orderBy []exp.OrderedExpression, columns []any, limit, offset int64) (ret []*T, err error) {
	return g.list(ctx, g.conn, where, orderBy, columns, limit, offset)
}

func (g *genericModelConn[T]) ListTx(ctx context.Context, sess sqlx.Session, where []goqu.Expression, orderBy []exp.OrderedExpression, columns []any, limit, offset int64) (ret []*T, err error) {
	return g.list(ctx, sess, where, orderBy, columns, limit, offset)
}

func (g *genericModelConn[T]) update(ctx context.Context, conn sqlx.Session, newValue goqu.Record, where []goqu.Expression) (sql.Result, error) {
	query, args, err := g.GenUpdateSQL(newValue, where)
	if err != nil {
		return nil, errors.Wrapf(err, "update prepare error")
	}
	return conn.ExecCtx(ctx, query, args...)
}

func (g *genericModelConn[T]) Update(ctx context.Context, newValue goqu.Record, where []goqu.Expression) (sql.Result, error) {
	return g.update(ctx, g.conn, newValue, where)
}

func (g *genericModelConn[T]) UpdateTx(ctx context.Context, sess sqlx.Session, newValue goqu.Record, where []goqu.Expression) (sql.Result, error) {
	return g.update(ctx, sess, newValue, where)
}

func (g *genericModelConn[T]) Sum(ctx context.Context, col exp.IdentifierExpression, where []goqu.Expression) (sum int64, err error) {
	return g.sum(ctx, g.conn, col, where)
}

func (g *genericModelConn[T]) SumTx(ctx context.Context, sess sqlx.Session, col exp.IdentifierExpression, where []goqu.Expression) (sum int64, err error) {
	return g.sum(ctx, sess, col, where)
}

func (g *genericModelConn[T]) GenUpdateSQL(newValue goqu.Record, where []goqu.Expression) (sql string, params []interface{}, err error) {
	return g.dialect.Update(g.tableName).Set(newValue).Where(where...).Prepared(true).ToSQL()
}

func (g *genericModelConn[T]) GenInsertSQL(model *T, onConflictDoNothing bool) (sql string, params []interface{}, err error) {
	fieldNames := builder.RawFieldNames(model)
	fieldNames = stringx.Remove(fieldNames, "`id`", "`create_time`", "`update_at`", "`updated_at`", "`update_time`", "`create_at`", "`created_at`")

	cols := []string{}
	lo.ForEach(fieldNames, func(x string, _ int) {
		cols = append(cols, strings.Trim(x, "`"))
	})

	vals, err := GetValsByTags(model, "db", cols)
	if err != nil {
		return "", nil, errors.Wrapf(err, "get valus by tags error")
	}

	ds := g.dialect.Insert(g.tableName)
	if onConflictDoNothing {
		ds = ds.OnConflict(goqu.DoUpdate("create_time", goqu.C("create_time").Set(goqu.C("create_time"))))
	}

	return ds.Cols(ConvertToInterfaceSlice(cols)...).Vals(goqu.Vals(vals)).Prepared(true).ToSQL()
}

func (g *genericModelConn[T]) GenSelectForUpdateSQL(where []goqu.Expression, Block bool) (sql string, params []interface{}, err error) {
	wait := exp.Wait
	if !Block {
		wait = exp.NoWait
	}
	return g.dialect.From(g.tableName).ForUpdate(wait).Where(where...).Prepared(true).ToSQL()
}

func (g *genericModelConn[T]) SelectOneForUpdateTx(ctx context.Context, sess sqlx.Session, where []goqu.Expression, block bool) (*T, error) {
	sql, args, err := g.GenSelectForUpdateSQL(where, block)
	if err != nil {
		return nil, errors.Wrapf(err, "gen select one for update sql error")
	}

	var resp T
	err = sess.QueryRowCtx(ctx, &resp, sql, args...)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// ---------------------------
func ConvertToInterfaceSlice[T any](input []T) []any {
	ret := []any{}
	lo.ForEach(input, func(x T, _ int) {
		ret = append(ret, x)
	})
	return ret
}

var ErrExpectStruct = fmt.Errorf("expect struct")
var ErrFieldCount = fmt.Errorf("field count error")

func GetValsByTags(stru any, tagName string, tagValues []string) ([]any, error) {
	typ := reflect.TypeOf(stru)
	val := reflect.ValueOf(stru)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, ErrExpectStruct
	}

	valMap := map[string]any{}

	for i := 0; i < typ.NumField(); i++ {
		tagVal := typ.Field(i).Tag.Get(tagName)
		if lo.Contains(tagValues, tagVal) {
			valMap[tagVal] = val.Field(i).Interface()
		}
	}

	ret := []any{}
	for _, t := range tagValues {
		if v, ok := valMap[t]; ok {
			ret = append(ret, v)
		} else {
			return nil, ErrFieldCount
		}
	}
	return ret, nil
}