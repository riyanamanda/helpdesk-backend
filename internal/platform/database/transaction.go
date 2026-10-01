package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Tx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	Commit() error
	Rollback() error
}

type TxManager struct {
	db *sqlx.DB
}

type transaction struct {
	tx *sqlx.Tx
}

func NewTxManager(db *sqlx.DB) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (m *TxManager) Begin(ctx context.Context) (Tx, error) {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &transaction{
		tx: tx,
	}, nil
}

func (t *transaction) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	return t.tx.QueryRowxContext(ctx, query, args...)
}

func (t *transaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}
