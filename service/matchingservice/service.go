package matchingservice

import (
	"GameApp/entity"
	"GameApp/param"
	"GameApp/pkg/richerror"
	"time"
)

type Repo interface {
	AddToWaitingList(UserID uint, category entity.Category) error
}
type Config struct {
	WaitingTimeOut time.Duration `koanf:"waiting_timeout"`
}

type Service struct {
	config Config
	repo   Repo
}

func New(config Config, repo Repo) Service {
	return Service{config: config, repo: repo}
}

func (s Service) AddToWaitingList(req param.AddToWaitingListRequest) (
	param.AddToWaitingResponse, error) {
	//	Get validated Data from Request
	const op = richerror.Op("matchinghandler.AddToWaitingList")

	// add user to the waiting list for the given category if not exist
	// also we can update the waiting timestamp
	err := s.repo.AddToWaitingList(req.UserID, req.Category)
	if err != nil {
		return param.AddToWaitingResponse{}, richerror.New(op).WithKind(richerror.KindUnexpected)
	}
	return param.AddToWaitingResponse{Timeout: s.config.WaitingTimeOut}, nil

}
