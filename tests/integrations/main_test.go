package integrations

import (
	"GameApp/adaptor/redis"
	"GameApp/pkg/testhelper"
	"fmt"
	"github.com/ory/dockertest/v3"
	"os"
	"testing"
)

var redisport string

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
	fmt.Println("redis port is:", redisRes)
	redisport = redisRes.GetPort("6379/tcp")
	exitCode := m.Run()
	os.Exit(exitCode)

}
