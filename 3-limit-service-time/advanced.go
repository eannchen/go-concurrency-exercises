package main

import (
	"sync"
	"time"
)

type entry struct {
	mu       sync.Mutex
	timeUsed int64
}

// AdvancedHandler implements the "10s max per user (accumulated)" logic.
type AdvancedHandler struct {
	mu        sync.Mutex
	userLimit map[int]*entry
}

// NewAdvancedHandler creates a new AdvancedHandler.
func NewAdvancedHandler() RequestHandler {
	return &AdvancedHandler{
		userLimit: make(map[int]*entry),
	}
}

// HandleRequest will contain the solution for the advanced level.
func (h *AdvancedHandler) HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	h.mu.Lock()
	e, ok := h.userLimit[u.ID]
	if !ok {
		e = &entry{}
		h.userLimit[u.ID] = e
	}

	// lock the user for user update
	e.mu.Lock()
	defer e.mu.Unlock()

	// `h.mu` is used for protecting map integrity,
	// releasing it to allow other users' request.
	h.mu.Unlock()

	remainingSecs := (10 - e.timeUsed) * int64(time.Second)
	if remainingSecs <= 0 {
		return false
	}

	timer := time.NewTimer(time.Duration(remainingSecs))
	defer func() { // the code explanation of this part is in `beginner.go`
		timer.Stop()
		select {
		case <-timer.C:
		default:
		}
	}()

	var start time.Time

	done := make(chan struct{})
	go func() {
		start = time.Now()
		process()
		close(done)
	}()

	select {
	case <-timer.C:
		e.timeUsed += int64(time.Since(start).Seconds())
		return false
	case <-done:
		e.timeUsed += int64(time.Since(start).Seconds())
		return true
	}
}
