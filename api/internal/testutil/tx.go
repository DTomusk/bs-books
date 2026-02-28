package testutil

import (
	"bs-books-api/internal/db"
	"context"
	"database/sql"
	"testing"
)

type TestTxRunner struct {
	tx *sql.Tx
}

func NewTestTxRunner(tx *sql.Tx) *TestTxRunner {
	return &TestTxRunner{tx: tx}
}

func (r *TestTxRunner) WithTx(
	ctx context.Context,
	fn func(tx *sql.Tx) error,
) error {
	return fn(r.tx)
}

func WithTx(t *testing.T, fn func(tx *sql.Tx)) {
	t.Helper()
	db := GetTestDB()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin tx: %v", err)
	}

	defer tx.Rollback()

	fn(tx)
}

func (r *TestTxRunner) DB() db.DBTX {
	return r.tx
}
