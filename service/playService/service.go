package playService

import (
	"GameApp/contract/broker"
	"GameApp/entity"
	"GameApp/param"
	"context"
)

type PlayRepo interface {
	createGame(ctx context.Context, category entity.Category, questions []entity.Question) (entity.Game, error)
	createPlayer(ctx context.Context, userID entity.User, gameID uint) (int, error)
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

// CreateGame make a game func
func (s Service) CreateGame(ctx context.Context, req param.CreateGameRequest, questions []entity.Question) (entity.Game, error) {
	game, err := s.repo.createGame(ctx, req.Category, questions)
	if err != nil {
		return entity.Game{}, err
	}
	return game, nil

}

// CreatePlayer create player func
func (s Service) CreatePlayer(ctx context.Context, userID uint, gameID uint) (int, error) {
	playerID, err := s.repo.createPlayer(ctx, entity.User{
		ID: userID,
	}, gameID)
	if err != nil {
		return playerID, err
	}
	return playerID, nil
}

// CreatePlayerAnswer create player Answer func
func (s Service) CreatePlayerAnswer(playerID entity.User, questionID entity.Question, choice entity.PossibleAnswerChoice) error {
	err := s.repo.createPlayerAnswer(playerID, questionID, choice)
	if err != nil {
		return err
	}
	return nil

}
