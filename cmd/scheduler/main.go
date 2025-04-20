package main

import (
	"GameApp/conf"
	"GameApp/scheduler"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("start Scheduler Service ....")
	cfg := conf.Load()
	fmt.Println(cfg)
	// TODO add command for migrations to dont run automatically

	done := make(chan bool)
	go func() {
		sch := scheduler.New()
		sch.Start(done)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Stopping Scheduler Service ...")
	done <- true
	time.Sleep(5 * time.Second)

}
