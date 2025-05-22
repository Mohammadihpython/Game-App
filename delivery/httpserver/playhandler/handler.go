package playhandler

import (
	"GameApp/contract/broker"
	"GameApp/entity"
	"GameApp/param"
	"GameApp/pkg/protobufEncoder"
	"GameApp/pkg/richerror"
	"GameApp/service/playService"
	"GameApp/service/questionservice"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Config struct {
	questionLimit uint
}

type Handler struct {
	consumer    broker.Consumer
	QuestionSVC questionservice.Service
	PlayerSVC   playService.Service
	Config      Config
}

func (h Handler) GameHandler(c echo.Context) error {
	const OP = "PlayHandler.GameHandler"
	//get the user ids and game category from event
	msg := h.consumer.Consumer(entity.MatchingUsersMatchedEvent)
	matchedUsers := protobufEncoder.DecoderEvent(entity.MatchingUsersMatchedEvent, msg)
	//get questions by category
	category := matchedUsers.Category
	questions, err := h.QuestionSVC.GetQuestions(string(category), h.Config.questionLimit)
	if err != nil {
		return richerror.New(OP).
			WithWrappedError(err).
			WithMessage("can not get Questions").
			WithKind(richerror.KindInvalid)

	}
	ctx := c.Request().Context()

	// create a game  and player and send message for client and users
	game, err := h.PlayerSVC.CreateGame(ctx, param.CreateGameRequest{Category: category}, questions)
	if err != nil {
		return richerror.New(OP).WithWrappedError(err)
	}
	var PlayersID []int
	for _, userid := range matchedUsers.UserIDs {
		palyerID, err := h.PlayerSVC.CreatePlayer(ctx, userid, game.ID)
		if err != nil {
			return richerror.New(OP).
				WithWrappedError(err).
				WithMessage("can not create player").
				WithKind(richerror.KindInvalid)

		}
		PlayersID = append(PlayersID, palyerID)

	}

	// send questions to clients
	// calculate user answers one by one and send response and update leader borde
	// send the another user data to user like how many question solve and
	var (
		upgrader = websocket.Upgrader{}
	)
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write

		err = ws.WriteMessage(1, []byte(game.Category))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
