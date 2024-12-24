package database

import (
	"fmt"
	"github.com/akmalfairuz/df-practice/practice/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error ping db: %w", err)
	}

	return db, nil
}

var globalDB *sqlx.DB

func Get() *sqlx.DB {
	return globalDB
}

func init() {
	if globalDB != nil {
		panic("database: globalDB is already initialized")
	}

	db, err := connect(config.Get().DSN())
	if err != nil {
		panic(fmt.Errorf("failed to connect to db: %w", err))
	}

	globalDB = db
}
