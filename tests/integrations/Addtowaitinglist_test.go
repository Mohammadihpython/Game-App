package integrations

import (
	presenceClient "GameApp/adaptor/presence"
	"GameApp/adaptor/redis"
	"GameApp/contract/goproto/presence"
	"GameApp/delivery/grpcserver/presenceserver"
	"GameApp/delivery/httpserver/matchinghandler"
	"GameApp/entity"
	"GameApp/param"
	"GameApp/repository/mysql"
	"GameApp/repository/mysql/mysqluser"
	redispresence "GameApp/repository/redis/presence"
	"GameApp/repository/redis/redismatching"
	"GameApp/service/authservice"
	"GameApp/service/matchingservice"
	"GameApp/service/presenceservice"
	"GameApp/service/userservice"
	"GameApp/validator/matchingsvalidator"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddToWaitingList(t *testing.T) {
	/*
		Start an In-Memory gRPC Server for Testing
	*/
	redisAdaptor := redis.New(cfg.Redis)
	presenceRepo := redispresence.New(redisAdaptor)
	svc := presenceservice.New(cfg.Presence, presenceRepo)
	presenceServer := presenceserver.New(svc)
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	presence.RegisterPresenceServiceServer(s, &presenceServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// start a user service and create a user`
	fmt.Println("mysql port", cfg.Mysql.Port)
	mysqlRepo := mysql.New(mysql.Config{
		Host:     cfg.Mysql.Host,
		Port:     cfg.Mysql.Port,
		Username: cfg.Mysql.Username,
		Password: cfg.Mysql.Password,
		DBName:   cfg.Mysql.DBName,
	})
	authSvc := authservice.New(cfg.Auth)

	userMysql := mysqluser.New(mysqlRepo)
	userSvc := userservice.New(userMysql, authSvc)
	createdUser, err := userSvc.Register(param.RegisterRequest{
		Name:        "Hamed",
		PhoneNumber: "09933642792",
		Password:    "1234",
	})
	if err != nil {
		t.Error(err)
	}

	// Start echo server
	e := echo.New()
	// create a request and body and get response
	reqBody := param.AddToWaitingListRequest{
		UserID:   createdUser.User.ID,
		Category: entity.FootballCategory,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/match/wait", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	// create matching service and send the request to its handler

	matchingRepo := redismatching.New(cfg.RedisMatching, redisAdaptor)
	matchingValidator := matchingsvalidator.New()

	// create in memory grpc server
	// DialOption configures how we set up the connection.
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	presenceC := presenceClient.New(
		"bufnet", // dummy address
		grpc.WithContextDialer(bufDialer),
	)

	matchingSVC := matchingservice.New(cfg.MatchingService, matchingRepo, presenceC, redisAdaptor)

	// Call handler
	handler := matchinghandler.New(cfg.Auth, authSvc, matchingSVC, matchingValidator, presenceC)
	err = handler.AddToWaitingList(ctx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

}
