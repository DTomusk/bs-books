package db

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type TxRunner interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
	DB() DBTX
}

type DBTxRunner struct {
	db *sql.DB
}

func NewDBTxRunner(db *sql.DB) *DBTxRunner {
	return &DBTxRunner{db: db}
}

func (r *DBTxRunner) WithTx(
	ctx context.Context,
	fn func(tx *sql.Tx) error,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *DBTxRunner) DB() DBTX {
	return r.db
}
