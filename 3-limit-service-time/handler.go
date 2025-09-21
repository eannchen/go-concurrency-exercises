package main

// User defines the UserModel.
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// RequestHandler defines the interface for handling user requests with time limits.
type RequestHandler interface {
	HandleRequest(process func(), u *User) bool
}
