package scheduler

import (
	"GameApp/param"
	"GameApp/service/matchingservice"
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	//get lock
	_, err := s.matchSVC.MatchWaitedUsers(ctx, param.MatchWaiteUserRequest{})
	if err != nil {
		//	TODO -log err
		//	TODO -update metrics
		fmt.Println("match waited users error:", err)
	}
	// free lock
}
