package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func DBConnect() {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")

	if dbHost == "" || dbName == "" || dbUser == "" || os.Getenv("STORAGE_PATH") == "" {
		panic("Missing database environment variabled!")
	}

	// TODO: retry connection or only connect at the first moment that we need the db
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?dial_timeout=10s&sslmode=disable", dbUser, dbPasswd, dbHost, dbName)

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))
	db = bun.NewDB(sqlDB, pgdialect.New())
}

func GetDB() *bun.DB {
	return db
}
