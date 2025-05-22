package entity

import "time"

type Game struct {
	ID          uint
	Category    string
	QuestionIDs []Question
	StartTime   time.Time
	WinnerID    uint
}

type Player struct {
	ID      uint
	UserID  uint
	GameID  uint
	Score   uint
	Answers []PlayerAnswer
}
type PlayerAnswer struct {
	ID         uint
	PlayerID   uint
	QuestionID uint
	Choice     PossibleAnswerChoice
}
