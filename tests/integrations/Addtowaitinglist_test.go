package integrations

import (
	"GameApp/param"
	"github.com/labstack/echo/v4"
	"testing"
)

func TestAddToWaitingList(t *testing.T) {
	e := echo.New()
	reqBody := param.AddToWaitingListRequest{
		UserID:   "",
		Category: "",
	}

}
