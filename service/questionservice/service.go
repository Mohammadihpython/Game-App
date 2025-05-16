package questionservice

import "GameApp/entity"

type Config struct {
}

type GameRepository interface {
	getQuestion(category string, questionLimit uint) ([]entity.Question, error)
}

type Service struct {
	config Config
	repo   GameRepository
}

func NewService(config Config) *Service {
	return &Service{config: config}
}

func (s Service) CreateQuestion() {}
func (s Service) GetQuestion()    {}
func (s Service) UpdateQuestion() {}
func (s Service) DeleteQuestion() {}
func (s Service) GetAnswers()     {}
func (s Service) GetAnswer()      {}

func (s Service) GetQuestions(category string, questionsLimit uint) ([]entity.Question, error) {
	questions, err := s.repo.getQuestion(category, questionsLimit)
	if err != nil {
		return nil, err
	}
	return questions, nil

}
