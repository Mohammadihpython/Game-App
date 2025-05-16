package integrations

import (
	"GameApp/conf"
	"github.com/labstack/echo/v4"
	"testing"
)

func TestAddToWaitingList(t *testing.T) {
	cfg := conf.Load()

	e := echo.New()

	println(e)

}
