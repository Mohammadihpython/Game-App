package integrations

import (
	"GameApp/adaptor/redis"
	"GameApp/conf"
	"GameApp/pkg/testhelper"
	"GameApp/repository/migrator"
	"GameApp/repository/mysql"
	"context"
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"google.golang.org/grpc/test/bufconn"
	"os"
	"testing"
	"time"
)

var lis *bufconn.Listener
var RedisPort string
var mysqlRepo *sql.DB
var cfg = conf.Load()

const bufSize = 1024 * 1024

func TestMain(m *testing.M) {

	pool := testhelper.StartDockerPool()
	fmt.Println("start creating services")

	// set up redis container for test
	redisRes := testhelper.StartDockerInstance(pool, "redis/redis-stack-server", "latest",
		func(res *dockertest.Resource) error {
			RedisPort = res.GetPort("6379/tcp")
			fmt.Println("redis port is:", RedisPort)

			redisAdaptor := redis.New(redis.Config{
				Host: "localhost",
				Port: RedisPort,
			})
			err := redisAdaptor.Client().Ping(context.Background()).Err()
			if err != nil {
				return err
			}

			return nil
		})
	mysqlRes := testhelper.StartDockerInstance(pool, "mysql", "8.0", func(res *dockertest.Resource) error {
		MysqlPort := res.GetPort("3306/tcp")
		fmt.Println("MySQL port is:", MysqlPort)
		cfg.Mysql.Port = MysqlPort

		// Retry until MySQL is ready to accept connections
		err := pool.Retry(func() error {
			mysqlAdaptor := mysql.New(mysql.Config{
				Host:     cfg.Mysql.Host, // Use localhost for Docker-exposed ports
				Port:     MysqlPort,
				Username: cfg.Mysql.Username,
				Password: cfg.Mysql.Password,
				DBName:   cfg.Mysql.DBName,
			})
			db := mysqlAdaptor.Conn()
			pingErr := db.Ping()
			if pingErr != nil {
				fmt.Println("Waiting for MySQL...", pingErr)
			}
			return pingErr
		})

		if err != nil {
			fmt.Println("MySQL not ready in time:", err)
			return err
		}

		return nil
	}, "MYSQL_ROOT_PASSWORD="+cfg.Mysql.Password,
		"MYSQL_DATABASE="+cfg.Mysql.DBName,
		"MYSQL_USER="+cfg.Mysql.Username,
		"MYSQL_PASSWORD="+cfg.Mysql.Password)
	// اعمال مایگریشن ها بر روی دیتابیس
	migrator := migrator.New(mysql.Config{
		Host:     cfg.Mysql.Host,
		Port:     cfg.Mysql.Port,
		Username: cfg.Mysql.Username,
		Password: cfg.Mysql.Password,
		DBName:   cfg.Mysql.DBName,
	}, "../../repository/mysql/migrations")
	time.Sleep(5 * time.Second)
	migrator.Up()
	//err := testhelper.RunMigrations(mysqlRepo, "../../repository/mysql/migrations/")
	//if err != nil {
	//	fmt.Println(err)
	//}
	fmt.Println(mysqlRes)
	fmt.Println("redis port is:", redisRes)
	RedisPort = redisRes.GetPort("6379/tcp")
	exitCode := m.Run()
	os.Exit(exitCode)

}
