package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgresql driver, don't remove
)

var db *sqlx.DB

//InitDB establishes connection to database and saves its handler into db *sqlx.DB
func InitDB(connection string) {
	db = sqlx.MustConnect("postgres", connection)
}

//GetDB returns database handler
func GetDB() *sqlx.DB {
	return db
}

//utility functions

//truncate truncates string to n runes
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}
