package integrations

import (
	"GameApp/adaptor/redis"
	"GameApp/pkg/testhelper"
	"GameApp/repository/mysql"
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"os"
	"testing"
)

var RedisPort string
var MysqlPort string
var mysqlRepo *sql.DB

func TestMain(m *testing.M) {
	pool := testhelper.StartDockerPool()

	fmt.Println("start creating services")

	// set up redis container for test
	redisRes := testhelper.StartDockerInstance(pool, "redis/redis-stack-server", "latest",
		func(res *dockertest.Resource) error {
			port := res.GetPort("6379/tcp")
			fmt.Println("redis port is:", port)

			redisAdaptor := redis.New(redis.Config{
				Host: "localhost",
				Port: port,
			})
			redisAdaptor.Client()
			return nil
		})
	mysqlRes := testhelper.StartDockerInstance(pool, "mysql", "8.0", func(res *dockertest.Resource) error {
		port := res.GetPort("3306/tcp")
		fmt.Println("mysql port is:", port)
		mysqlAdaptor := mysql.New(mysql.Config{
			Host:     "localhost",
			Port:     port,
			Username: "hamed",
			Password: "1234",
			DBName:   "test_db",
		})
		mysqlRepo = mysqlAdaptor.Conn()
		return nil
	})
	// اعمال مایگریشن ها بر روی دیتابیس
	err := testhelper.RunMigrations(mysqlRepo, ".../repository/migrations/")
	if err != nil {
		fmt.Println(err)
	}
	MysqlPort = mysqlRes.GetPort("3307/tcp")
	fmt.Println("redis port is:", redisRes)
	RedisPort = redisRes.GetPort("6379/tcp")
	exitCode := m.Run()
	os.Exit(exitCode)

}
