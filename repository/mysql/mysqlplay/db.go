package mysqlplay

import "GameApp/repository/mysql"

type DB struct {
	conn *mysql.MYSQL
}

func New(conn *mysql.MYSQL) *DB {
	return &DB{conn}
}
