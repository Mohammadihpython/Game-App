package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

type MYSQL struct {
	config Config
	db     *sql.DB
}

func New(cfg Config) *MYSQL {
	path := fmt.Sprintf("%s:%s@(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	fmt.Println(path)
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName))
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &MYSQL{config: cfg, db: db}

}
