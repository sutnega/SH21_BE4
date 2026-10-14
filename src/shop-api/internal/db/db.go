package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func Connect(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(20)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}
