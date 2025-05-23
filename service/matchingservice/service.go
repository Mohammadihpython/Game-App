package matchingservice

import (
	"GameApp/contract/broker"
	"GameApp/entity"
	"GameApp/param"
	"GameApp/pkg/protobufEncoder"
	"GameApp/pkg/richerror"
	"GameApp/pkg/timestamp"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"sync"
	"time"
)

type Repo interface {
	GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
	AddToWaitingList(UserID uint, category entity.Category) error
	RemoveUsersFromWaitingList(userIDs []uint, category entity.Category)
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
	pub            broker.Publisher
}

func New(config Config, repo Repo, presenceClient PresenceClient, pub broker.Publisher) Service {
	return Service{config: config, repo: repo, presenceClient: presenceClient, pub: pub}
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

func (s Service) MatchWaitedUsers(ctx context.Context, _ param.MatchWaiteUserRequest) (param.MatchWaiteUserResponse, error) {
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
	fmt.Println("start job")
	defer wg.Done()
	list, err := s.repo.GetWaitingListByCategory(ctx, category)
	fmt.Println("the list of waited uses", list, err)

	if err != nil {
		log.Errorf("GetWaitingListByCategory err:%v\n", err)
		return
	}
	userIDs := make([]uint, 0)
	for _, l := range list {
		userIDs = append(userIDs, l.UserID)
	}
	fmt.Println("userIDS : ", userIDs)
	if len(userIDs) < 2 {
		return
	}
	presenceList, err := s.presenceClient.GetPresence(ctx, param.GetPresenceRequest{UserIDs: userIDs})
	if err != nil {

		fmt.Println("match err", err)
	}

	presenceUserIDs := make([]uint, 0, len(list))
	for _, l := range presenceList.Items {
		presenceUserIDs = append(presenceUserIDs, l.UserID)
	}

	fmt.Println("presenceUserIDs : ", presenceUserIDs)
	var toBeRemovedUsers = make([]uint, 0)
	var finalList = make([]entity.WaitingMember, 0)
	for _, l := range list {
		//اینجا ما اولن زمان انلاین بودن کاربر رو چک میکنیم و
		// بعد اگر از ۲۰ ثانیه قبل تر بیشتر بود به final list اضاقه اش میکنیم
		LastOnlineTimestamp, ok := getPresenceItem(presenceList, l.UserID)
		if ok && LastOnlineTimestamp > timestamp.Add(-20*time.Second) &&
			l.Timestamp > timestamp.Add(-300*time.Second) {

			finalList = append(finalList, l)
		} else {
			toBeRemovedUsers = append(toBeRemovedUsers, l.UserID)

		}
	}
	go s.repo.RemoveUsersFromWaitingList(toBeRemovedUsers, category)
	matchedUsersTobeRemoved := make([]uint, 0)
	for i := 0; i < len(list)-1; i = i + 2 {
		matchedUsers := entity.MatchedUsers{
			Category: category,
			UserIDs:  []uint{finalList[i].UserID, finalList[i+1].UserID},
		}
		// publish a new event for mu
		payload := protobufEncoder.EncodeEvent(entity.MatchingUsersMatchedEvent, matchedUsers)
		fmt.Println("this is payload", payload)
		go s.pub.Publish(entity.MatchingUsersMatchedEvent, payload)

		matchedUsersTobeRemoved = append(matchedUsersTobeRemoved, userIDs...)
	}

	go s.repo.RemoveUsersFromWaitingList(matchedUsersTobeRemoved, category)
}

func getPresenceItem(presenceList param.GetPresenceResponse, userID uint) (int64, bool) {
	for _, item := range presenceList.Items {
		if item.UserID == userID {
			return item.Timestamp, true
		}
	}
	return 0, false
}
