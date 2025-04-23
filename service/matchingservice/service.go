package matchingservice

import (
	"GameApp/entity"
	"GameApp/param"
	"GameApp/pkg/richerror"
	"GameApp/pkg/timestamp"
	"context"
	"fmt"
	funk "github.com/thoas/go-funk"
	"sync"
	"time"
)

type Repo interface {
	GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
	AddToWaitingList(UserID uint, category entity.Category) error
}

// PresenceClient We use Grpc call to get presence
type PresenceClient interface {
	GetPresence(ctx context.Context, request param.GetPresenceRequest) (param.GetPresenceResponse, error)
}
type Config struct {
	WaitingTimeOut time.Duration `koanf:"waiting_timeout"`
}

type Service struct {
	config         Config
	repo           Repo
	presenceClient PresenceClient
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

func (s Service) MatchWaitedUsers(ctx context.Context, req param.MatchWaiteUserRequest) (param.MatchWaiteUserResponse, error) {
	const OP = richerror.Op("matchingservice.MatchWaitedUsers")
	wg := sync.WaitGroup{}
	for _, category := range entity.CategoryList() {
		wg.Add(1)
		go s.match(ctx, category, &wg)
	}
	wg.Wait()

	return param.MatchWaiteUserResponse{}, nil
	//create a new matched event(message) and publish it to broker

}

func (s Service) match(ctx context.Context, category entity.Category, wg *sync.WaitGroup) {
	// pre compute
	defer wg.Done()
	list, err := s.repo.GetWaitingListByCategory(ctx, category)
	if err != nil {
		return
	}
	userIDs := make([]uint, 0, len(list))
	for _, l := range list {
		userIDs = append(userIDs, l.UserID)
	}
	presenceList, err := s.presenceClient.GetPresence(ctx, param.GetPresenceRequest{UserIDs: userIDs})
	if err != nil {
		return
	}

	// Todo merge presenceList with list based on userID
	// also consider the presence timestamp of each user
	// and remove user from waiting list if the user time stamp older than time.Now(-20 second)
	//if t < timestamp.Add(-20*time.Second) {
	//		remove list[i].userID from waiting List
	//
	//}

	presenceUserIDs := make([]uint, 0, len(list))
	for _, l := range presenceList.Items {
		presenceUserIDs = append(presenceUserIDs, l.UserID)
	}
	var finalList = make([]entity.WaitingMember, 0)
	for _, l := range list {
		if funk.ContainsUInt(presenceUserIDs, l.UserID) && l.Timestamp < timestamp.Add(-20*time.Second) {
			finalList = append(finalList, l)
		} else {
			//	remove from list
		}
	}
	for i := 0; i < len(list)-1; i = i + 2 {

		mu := entity.MatchedUsers{
			Category: category,
			UserID:   []uint{finalList[i].UserID, finalList[i+1].UserID},
		}
		fmt.Println(mu)

		//	 publish a new event for mu
		//	 remove mu users from waiting list
	}
}
