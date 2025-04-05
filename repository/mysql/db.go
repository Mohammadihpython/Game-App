package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type MYSQL struct {
	db *sql.DB
}

func New() *MYSQL {

	db, err := sql.Open("mysql", "Hamed:hmah8013@(localhost:3308)/gameappDB")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &MYSQL{db}

}
