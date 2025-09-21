//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import "fmt"

func main() {
	handlers := map[string]RequestHandler{
		"Beginner": NewBeginnerHandler(),
		"Advanced": NewAdvancedHandler(),
	}

	for name, handler := range handlers {
		fmt.Printf("--- Running Mock Server with %s Handler ---\n", name)
		RunMockServer(handler)
		fmt.Println()
	}
}
