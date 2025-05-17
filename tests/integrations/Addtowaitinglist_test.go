package integrations

import (
	presenceClient "GameApp/adaptor/presence"
	"GameApp/adaptor/redis"
	"GameApp/conf"
	"GameApp/contract/goproto/presence"
	"GameApp/delivery/httpserver/matchinghandler"
	"GameApp/entity"
	"GameApp/param"
	"GameApp/repository/mysql"
	"GameApp/repository/mysql/mysqluser"
	"GameApp/repository/redis/redismatching"
	"GameApp/service/authservice"
	"GameApp/service/matchingservice"
	"GameApp/service/presenceservice"
	"GameApp/service/userservice"
	"GameApp/validator/matchingsvalidator"
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddToWaitingList(t *testing.T) {
	cfg := conf.Load()
	mysqlRepo := mysql.New(cfg.Mysql)
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
	e := echo.New()
	reqBody := param.AddToWaitingListRequest{
		UserID:   createdUser.User.ID,
		Category: entity.FootballCategory,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/match/wait", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	redisAdaptor := redis.New(cfg.Redis)
	matchingRepo := redismatching.New(cfg.RedisMatching, redisAdaptor)
	matchingValidator := matchingsvalidator.New()

	presenceSVC := presenceservice.New(cfg.Presence, mysqlRepo)
	matchingSVC := matchingservice.New(cfg.MatchingService, matchingRepo, presenceSClient, redisAdaptor)

	// Call handler
	handler := matchinghandler.New(cfg.Auth, authSvc, matchingSVC, matchingValidator, presenceSClient)
	err = handler.AddToWaitingList(ctx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

}
