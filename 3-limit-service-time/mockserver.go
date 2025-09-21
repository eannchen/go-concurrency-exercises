//////////////////////////////////////////////////////////////////////
//
// DO NOT EDIT THIS PART
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// RunMockServer simulates user interactions with the service.
// It now takes a handler to test a specific implementation.
func RunMockServer(handler RequestHandler) {
	var wg sync.WaitGroup
	u1 := User{ID: 1, IsPremium: false} // Free User
	u2 := User{ID: 2, IsPremium: true}  // Premium User

	requests := []struct {
		pid     int
		process func()
		user    *User
		delay   time.Duration
	}{
		{1, shortProcess, &u1, 1 * time.Second}, // u1 uses 6s. Total: 6s.
		{2, longProcess, &u2, 2 * time.Second},  // u2 is premium, should pass.
		{3, shortProcess, &u1, 1 * time.Second}, // u1 uses another 6s. Total: 12s. Should be killed by advanced handler.
		{4, longProcess, &u1, 0 * time.Second},  // u1 has no time left. Should be killed. Also >10s, killed by beginner.
		{5, shortProcess, &u2, 0 * time.Second}, // u2 is premium, should pass.
	}

	wg.Add(len(requests))

	for _, req := range requests {
		// create a new variable for the goroutine
		r := req
		go func() {
			defer wg.Done()
			createMockRequest(handler, r.pid, r.process, r.user)
		}()
		time.Sleep(r.delay)
	}

	wg.Wait()
}

func createMockRequest(handler RequestHandler, pid int, fn func(), u *User) {
	fmt.Println("UserID:", u.ID, "\tProcess", pid, "started.")
	res := handler.HandleRequest(fn, u)

	if res {
		fmt.Println("UserID:", u.ID, "\tProcess", pid, "done.")
	} else {
		fmt.Println("UserID:", u.ID, "\tProcess", pid, "killed. (No quota left)")
	}
}

func shortProcess() {
	time.Sleep(6 * time.Second)
}

func longProcess() {
	time.Sleep(11 * time.Second)
}
