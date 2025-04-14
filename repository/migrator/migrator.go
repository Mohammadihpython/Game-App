package migrator

import (
	"GameApp/repository/mysql"

	"database/sql"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	cfg        mysql.Config
	migrations *migrate.FileMigrationSource
}

func New(cfg mysql.Config) Migrator {

	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}
	return Migrator{cfg: cfg, migrations: migrations}
}

func (m Migrator) Up() {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", m.cfg.Username, m.cfg.Password, m.cfg.Host, m.cfg.Port, m.cfg.DBName))
	if err != nil {
		panic(fmt.Errorf("cant open mysql db %v", err))
	}
	n, err := migrate.Exec(db, "mysql", m.migrations, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("failed to execute migrations: %w", err))
	}
	fmt.Printf("%d applied migrations\n", n)

}
func (m Migrator) Down(cfg mysql.Config) {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName))
	if err != nil {
		panic(fmt.Errorf("cant open mysql db %v", err))
	}
	n, err := migrate.Exec(db, "mysql", m.migrations, migrate.Down)
	if err != nil {
		panic(fmt.Errorf("failed to RollBack migrations: %w", err))
	}
	fmt.Printf("%d Rollbacked migrations\n", n)

}
func (m Migrator) Status(cfg mysql.Config) {
	//	TODO Add status

}
