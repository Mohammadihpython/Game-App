package entity

type Game struct {
	ID          uint
	Category    uint
	QuestionIDs []uint
	PlayerIDs   []Player
	//WinnerID    uint
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
