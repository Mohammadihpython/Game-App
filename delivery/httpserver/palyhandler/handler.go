package palyhandler

import (
	"GameApp/contract/broker"
	"GameApp/entity"
	"GameApp/pkg/protobufEncoder"
	"GameApp/pkg/richerror"
	"GameApp/service/playService"
	"GameApp/service/questionservice"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"time"
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
	// get the user ids and game category from event
	msg := h.consumer.Consumer(entity.MatchingUsersMatchedEvent)
	matchedUsers := protobufEncoder.DecoderEvent(entity.MatchingUsersMatchedEvent, msg)
	// get questions by category
	category := matchedUsers.Category
	questions, err := h.QuestionSVC.GetQuestions(string(category), h.Config.questionLimit)
	if err != nil {
		richerror.New(OP).
			WithWrappedError(err).
			WithMessage("can not get Questions").
			WithKind(richerror.KindInvalid),
	)



	}

	// create an game  and player and send message for client and users
	// send questions to clients
	// calculate user answers one by one and send response and update leader borde
	// send the another user data to user like how many question solve and
	//
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {

			// Write

			err := websocket.Message.Send(ws, "Hello, Dear !")
			if err != nil {
				c.Logger().Error(err)
			}

			// Read
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func play(c echo.Context) error {
	userIDParam := c.QueryParam("user_id")
	var userID uint
	fmt.Sscanf(userIDParam, "%d", &userID)

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		userSockets[userID] = ws

		// 1. Get matched players and game info (mocked for now)
		otherUserID := uint(2) // You would fetch this dynamically
		game := Game{
			ID:          1,
			CategoryID:  3,
			PlayerIDs:   []uint{userID, otherUserID},
			QuestionIDs: []uint{101, 102, 103},
			StartTime:   time.Now(),
			Status:      "in_progress",
		}

		player := Player{
			ID:      1,
			UserID:  userID,
			GameID:  game.ID,
			Score:   0,
			Answers: []PlayerAnswer{},
		}

		// 2. Send start message
		startMsg := map[string]any{
			"event":     "game_start",
			"game_id":   game.ID,
			"questions": len(game.QuestionIDs),
		}
		websocket.JSON.Send(ws, startMsg)

		// 3. Iterate over questions one-by-one
		for i, qID := range game.QuestionIDs {
			// Send question
			question := map[string]any{
				"event":       "new_question",
				"question_id": qID,
				"question":    fmt.Sprintf("What is the answer to question %d?", qID),
				"choices":     []string{"A", "B", "C", "D"},
			}
			websocket.JSON.Send(ws, question)

			// Receive answer
			var answerMsg struct {
				Event      string `json:"event"`
				QuestionID uint   `json:"question_id"`
				Choice     string `json:"choice"`
			}
			if err := websocket.JSON.Receive(ws, &answerMsg); err != nil {
				c.Logger().Error("failed to receive answer:", err)
				return
			}

			// Mock answer checking
			correct := answerMsg.Choice == "A" // assume A is correct
			if correct {
				player.Score++
			}

			// Save the answer
			player.Answers = append(player.Answers, PlayerAnswer{
				QuestionID: qID,
				Choice:     answerMsg.Choice,
				Correct:    correct,
			})

			// 4. Send feedback
			websocket.JSON.Send(ws, map[string]any{
				"event":   "answer_result",
				"correct": correct,
				"score":   player.Score,
			})

			// 5. Update opponent’s view (if connected)
			if opponentWS, ok := userSockets[otherUserID]; ok {
				websocket.JSON.Send(opponentWS, map[string]any{
					"event":   "opponent_progress",
					"user_id": userID,
					"solved":  i + 1,
					"score":   player.Score,
				})
			}
		}

		// 6. Game over
		websocket.JSON.Send(ws, map[string]any{
			"event": "game_over",
			"score": player.Score,
		})

	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
