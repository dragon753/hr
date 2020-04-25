package hr

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestSimpleJob(t *testing.T) {
	manager := NewManager()

	jobName := "hello_job"
	cronExpress := "* * * * *"

	start := time.Now()
	manager.Register(jobName, cronExpress, func() {
		now := time.Now()
		log.Println("say hello once.", now)

		if now.Sub(start).Seconds() > 60 {
			t.Fatal("should less than 60")
		}
		if now.Second() != 0 {
			t.Fatal("should be zero")
		}
		manager.Quit()
	})
	manager.Start()
}

func TestPanicJob(t *testing.T) {
	wait := sync.WaitGroup{}
	wait.Add(1)

	errMsg := "some error"
	manager := NewManagerWithOutput(nil, func(job Job, err error) {
		if fmt.Sprintf("%v", err) != errMsg {
			t.Fatalf("should equal")
		}
		log.Println("panic job stops")
		wait.Done()
	})

	start := time.Now()
	manager.Register("panic_job", "* * * * *", func() {
		now := time.Now()
		if now.Sub(start).Seconds() > 60 {
			t.Fatal("should less than 60")
		}
		if now.Second() != 0 {
			t.Fatal("should be zero")
		}

		panic(fmt.Errorf(errMsg))
	})

	go func() {
		wait.Wait()
		manager.Quit()
	}()
	manager.Start()
}

func TestManyJobs(t *testing.T) {
	manager := NewManager()

	jobName := "hello_job"
	cronExpress := "* * * * *"

	manager.Register(jobName+" 0", cronExpress, func() {
		log.Println("say hello 0.", time.Now())
	})
	manager.Register(jobName+" 1", cronExpress, func() {
		log.Println("say hello 1.", time.Now())
	})
	manager.Register(jobName+" 2", cronExpress, func() {
		log.Println("say hello 2.", time.Now())
	})
	manager.Register(jobName+" last", cronExpress, func() {
		log.Println("say hello once.", time.Now())
		manager.Quit()
	})
	manager.Start()
}


func TestSecondlyJob(t *testing.T) {
	manager := NewManager()

	jobName := "secondly job"
	cronExpress := "* * * * * * *"

	start := time.Now()
	manager.Register(jobName, cronExpress, func() {
		manager.Quit()
		log.Println("say hello secondly.", time.Now())
		now := time.Now()
		if now.Sub(start).Seconds() >1 {
			t.Fatal("should less than 2")
		}
	})
	manager.Start()
}
