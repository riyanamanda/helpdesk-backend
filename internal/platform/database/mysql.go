package database

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func NewMySql(conn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", conn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

func RunMySqlMigrations(db *sqlx.DB) error {
	goose.SetDialect("mysql")
	return goose.Up(db.DB, "migrations")
}
