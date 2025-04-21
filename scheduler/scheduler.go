package scheduler

import (
	"GameApp/param"
	"GameApp/service/matchingservice"
	"fmt"
	"github.com/go-co-op/gocron"
	"sync"
	"time"
)

type Config struct {
	MatchWaitedUsersInternalInSeconds int `koanf:"match_waited_users_internal_in_seconds"`
}

type Scheduler struct {
	config   Config
	sch      *gocron.Scheduler
	matchSVC matchingservice.Service
}

func New(config Config, matchSVC matchingservice.Service) Scheduler {
	return Scheduler{config: config, sch: gocron.NewScheduler(time.UTC), matchSVC: matchSVC}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	s.sch.Every(s.config.MatchWaitedUsersInternalInSeconds).Second().Do(s.MatchWaitedUsers)
	s.sch.StartBlocking()

	<-done
	//wait to finish job
	fmt.Println("stop scheduler....")
	s.sch.Stop()

}

func (s Scheduler) MatchWaitedUsers() {
	resp, err := s.matchSVC.MatchWaitedUsers(param.MatchWaiteUserRequest{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), resp)
	//
}
