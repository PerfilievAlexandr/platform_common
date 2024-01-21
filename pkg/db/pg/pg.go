package pg

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"platform_common/pkg/db"
	"platform_common/pkg/db/transaction"
)

type pg struct {
	dbc *pgxpool.Pool
}

func New(ctx context.Context, connectStr string) (db.Client, error) {
	dbc, err := pgxpool.Connect(ctx, connectStr)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &pg{dbc}, nil
}

func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	row, err := p.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	row, err := p.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, row)
}

func (p *pg) ExecContext(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	tx, ok := ctx.Value(transaction.TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, query, args...)
	}

	return p.dbc.Exec(ctx, query, args...)
}

func (p *pg) QueryContext(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	tx, ok := ctx.Value(transaction.TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, query, args...)
	}

	return p.dbc.Query(ctx, query, args...)
}

func (p *pg) QueryRowContext(ctx context.Context, query string, args ...interface{}) pgx.Row {
	tx, ok := ctx.Value(transaction.TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, query, args...)
	}

	return p.dbc.QueryRow(ctx, query, args...)
}

func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

func (p *pg) Close() {
	p.dbc.Close()
}
