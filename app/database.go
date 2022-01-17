package app

import (
	"azmi17/go-rest-api/helper"
	"database/sql"
	"time"
)

func NewDB() *sql.DB {

	// Db Config
	db, err := sql.Open("mysql", "root@tcp(localhost:3317)/go_rest_api")
	helper.PanicIfError(err)

	// Db Pooling
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
