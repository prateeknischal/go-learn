package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// ServiceResponse
type ServiceResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Load - load response struct
type Load struct {
	ServiceResponse ServiceResponse
	Load            int32 `json:"load"`
}

type Job struct {
	Id    int
	Retry int
}

func main() {
	var wg sync.WaitGroup

	client := http.Client{
		Timeout: time.Duration(3 * time.Second),
	}

	req := make(chan Job, 100)
	ack := make(chan bool, 100)

	status := make(map[bool]int)
	// var for atomic int ops
	var count int32 = 0

	go func() {
		wg.Add(1)
		for i := 0; i < 100; i++ {
			// push jobs in the channel
			req <- Job{Id: i, Retry: 0}
			atomic.AddInt32(&count, 1)
		}
		fmt.Println("****** All Jobs pushed ******", count)

		for count != 0 {
			status[<-ack]++
			atomic.AddInt32(&count, -1)
		}

		close(ack)
		defer func() { close(req); wg.Done() }()
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			ok := true
			for ok {
				t, ok := <-req
				if !ok {
					break
				}
				resp, err := client.Get("http://localhost:8081/status")
				if err != nil {
					fmt.Println("Job ", t.Id, "failed, Err: ", err)
					t.Retry += 1
					if t.Retry < 3 {
						fmt.Println("Pushing Job ", t.Id, " with retry ", t.Retry)
						req <- t
					} else {
						fmt.Println("Job ", t.Id, " giving up!")
						ack <- false
					}
				} else {
					fmt.Println("Job ", t.Id, "Success, Status: ", resp.StatusCode, ", retry : ", t.Retry, " Count ", count)
					ack <- true
				}
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Status : Success :", status[true], ", Failed : ", status[false])
}
