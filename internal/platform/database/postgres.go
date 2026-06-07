package database

import (
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func NewPostgres(conn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		slog.Error("error connect to database", "error", err)
		panic(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)

	return db
}

func RunMigrations(db *sqlx.DB) error {
	goose.SetDialect("postgres")
	return goose.Up(db.DB, "migrations")
}
