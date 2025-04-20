package scheduler

import (
	"GameApp/param"
	"GameApp/service/matchingservice"
	"fmt"
	"github.com/go-co-op/gocron"
	"sync"
	"time"
)

type Scheduler struct {
	sch      *gocron.Scheduler
	matchSVC matchingservice.Service
}

func New(matchSVC matchingservice.Service) Scheduler {
	return Scheduler{sch: gocron.NewScheduler(time.UTC), matchSVC: matchSVC}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	s.sch.Every(5).Second().Do(s.MatchWaitedUsers)
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
	fmt.Println(resp)
	//
}
