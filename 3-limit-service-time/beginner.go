package main

import (
	"time"
)

// BeginnerHandler implements the "10s max per request" logic.
type BeginnerHandler struct{}

// NewBeginnerHandler creates a new BeginnerHandler.
func NewBeginnerHandler() RequestHandler {
	return &BeginnerHandler{}
}

// HandleRequest will contain the solution for the beginner level.
func (h *BeginnerHandler) HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	done := make(chan struct{})
	timer := time.NewTimer(10 * time.Second)
	defer func() {
		// If true, it means we stop the timer proactively,
		// the `timer.C` will have no value need to be exhausted.
		if timer.Stop() {
			return
		}

		// If false, it means the timer is stopped already,
		// the following code ensures exhausting the `timer.C`.
		select {
		case <-timer.C:
		default:
		}

		// No need to close `done` proactively,
		// GC will clean up when the goroutine exits,
		// keep it will have the risk of causing `panic: close of closed channel`.
		// close(done)
	}()

	go func() {
		process()
		close(done)
	}()

	select {
	case <-timer.C:
		return false
	case <-done:
		return true
	}
}
