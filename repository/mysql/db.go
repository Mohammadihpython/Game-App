package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type Config struct {
	Host     string `koanf:"host"`
	Port     string `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	DBName   string `koanf:"dbname"`
}

type MYSQL struct {
	config Config
	db     *sql.DB
}

func (m MYSQL) Conn() *sql.DB {
	return m.db
}

func New(cfg Config) *MYSQL {
	port, _ := strconv.Atoi(cfg.Port)
	add := fmt.Sprintf("%s:%s@(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, port, cfg.DBName)
	fmt.Println(add)
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, port, cfg.DBName))
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &MYSQL{config: cfg, db: db}

}
