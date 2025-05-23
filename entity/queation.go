package entity

type Question struct {
	ID              uint
	Text            string
	PossibleAnswers []PossibleAnswer
	CorrectAnswerID uint
	Difficulty      QuestionDifficulty
	Category        uint
}

type PossibleAnswer struct {
	ID     int
	Text   string
	choice PossibleAnswerChoice
}

type PossibleAnswerChoice uint8

func (p PossibleAnswerChoice) IsValid() bool {
	if p >= PossibleAnswerA && p <= PossibleAnswerD {
		return true
	}
	return false
}

// ENUM
const (
	PossibleAnswerA PossibleAnswerChoice = iota + 1
	PossibleAnswerB
	PossibleAnswerC
	PossibleAnswerD
)

type QuestionDifficulty uint8

const (
	QuestionDifficultyEasy QuestionDifficulty = iota + 1
	QuestionDifficultyMedium
	QuestionDifficultyHard
)

func (q QuestionDifficulty) IsValid() bool {
	if q >= QuestionDifficultyEasy && q <= QuestionDifficultyHard {
		return true
	}
	return false
}
