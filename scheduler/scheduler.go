package scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct {
}

func New() Scheduler {
	return Scheduler{}
}

func (s Scheduler) Start(done <-chan bool) {
	for {
		select {
		case <-done:
			fmt.Println("Scheduler stopped")
			return

		default:
			now := time.Now()
			fmt.Println(now)
			time.Sleep(3 * time.Second)

		}

	}

}
