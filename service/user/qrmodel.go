package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	qrFieldNames          = builder.RawFieldNames(&Qr{})
	qrRows                = strings.Join(qrFieldNames, ",")
	qrRowsExpectAutoSet   = strings.Join(stringx.Remove(qrFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	qrRowsWithPlaceHolder = strings.Join(stringx.Remove(qrFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheQrIdPrefix = "cache:qr:id:"
)

type (
	QrModel interface {
		Insert(ctx context.Context, data *Qr) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Qr, error)
		Update(ctx context.Context, data *Qr) error
		Delete(ctx context.Context, id int64) error
	}

	defaultQrModel struct {
		sqlc.CachedConn
		table string
	}

	Qr struct {
		Id         int64     `db:"id"`
		QrCode     string    `db:"qr_code"`    // 验证码
		AuthCount  int64     `db:"auth_count"` // 验证次数
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewQrModel(conn sqlx.SqlConn, c cache.CacheConf) QrModel {
	return &defaultQrModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`qr`",
	}
}

func (m *defaultQrModel) Insert(ctx context.Context, data *Qr) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, qrRowsExpectAutoSet)
	ret, err := m.ExecNoCacheCtx(ctx, query, data.QrCode, data.AuthCount)

	return ret, err
}

func (m *defaultQrModel) FindOne(ctx context.Context, id int64) (*Qr, error) {
	qrIdKey := fmt.Sprintf("%s%v", cacheQrIdPrefix, id)
	var resp Qr
	err := m.QueryRowCtx(ctx, &resp, qrIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", qrRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultQrModel) Update(ctx context.Context, data *Qr) error {
	qrIdKey := fmt.Sprintf("%s%v", cacheQrIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, qrRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.QrCode, data.AuthCount, data.Id)
	}, qrIdKey)
	return err
}

func (m *defaultQrModel) Delete(ctx context.Context, id int64) error {
	qrIdKey := fmt.Sprintf("%s%v", cacheQrIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, qrIdKey)
	return err
}

func (m *defaultQrModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheQrIdPrefix, primary)
}

func (m *defaultQrModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", qrRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}
