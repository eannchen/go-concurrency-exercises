//////////////////////////////////////////////////////////////////////
//
// DO NOT EDIT THIS PART
// Your task is to edit `main.go`
//

package main

import (
	"sync/atomic"
	"time"
)

// MockDB used to simulate a database model
type MockDB struct {
	// A counter to keep track of how many times the database's Get method has been called.
	Calls int32
}

// Get only returns the key, as this is only for demonstration purposes
func (db *MockDB) Get(key string) (string, error) {
	d, _ := time.ParseDuration("20ms")
	time.Sleep(d)
	// It uses atomic.AddInt32 instead of a simple db.Calls++ because the mock server calls this method from many goroutines at once.
	// The atomic operation guarantees that the counter is incremented correctly without causing a data race.
	atomic.AddInt32(&db.Calls, 1)
	return key, nil
}

// GetMockDB returns an instance of MockDB
func GetMockDB() *MockDB {
	return &MockDB{
		Calls: 0,
	}
}
