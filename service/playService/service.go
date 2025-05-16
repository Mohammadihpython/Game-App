package playService

import (
	"GameApp/contract/broker"
	"GameApp/entity"
)

type PlayRepo interface {
	createGame(category entity.Category, questionIDs []uint, playerIDs []uint) (entity.Game, error)
	createPlayer(userID entity.User, gameID entity.Game, score uint, AnswerIDs []entity.PlayerAnswer) (int, error)
	createPlayerAnswer(playerID entity.User, questionID entity.Question, choice entity.PossibleAnswerChoice) error
}

type Service struct {
	consumer broker.Consumer
	repo     PlayRepo
}

func NewPlayService(consumer broker.Consumer, player PlayRepo) Service {
	return Service{
		consumer: consumer,
		repo:     player,
	}
}

// implement create game

// make a game func
func (s Service) createGame(category entity.Category, questionIDs []uint, playerIDs []uint) (entity.Game, error) {
	gameID, err := s.repo.createGame(category, questionIDs, playerIDs)
	if err != nil {
		return entity.Game{}, err
	}
	return gameID, nil

}

// create player func
func (s Service) createPlayer(userID entity.User, gameID entity.Game, score uint, AnswerIDs []entity.PlayerAnswer) error {
	playerID, err := s.repo.createPlayer(userID, gameID, score, AnswerIDs)
	if err != nil {
		return err
	}
	return nil
}

// create player Answer func
func (s Service) createPlayerAnswer(playerID entity.User, questionID entity.Question, choice entity.PossibleAnswerChoice) error {
	err := s.repo.createPlayerAnswer(playerID, questionID, choice)
	if err != nil {
		return err
	}
	return nil

}
