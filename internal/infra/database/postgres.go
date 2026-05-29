package database

import (
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
