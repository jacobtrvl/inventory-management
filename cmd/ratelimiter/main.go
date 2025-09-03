// Copyright 2025 Jacob Philip. All rights reserved.
package main

import (
	"fmt"

	"github.com/jacobtrvl/inventory-management/pkg/ratelimiter"
)

func main() {
	rl := ratelimiter.NewRateLimiter(5)

	//Block the main loop with parallel requests
	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println("Requesting operation goroutine", i+1)
			rl.Allow()
			fmt.Println("Allowed operation goroutine", i+1)

		}(i)
	}
	for i := 0; i < 10; i++ {
		fmt.Println("Requesting operation", i+1)
		rl.Allow()
		fmt.Println("Allowed operation", i+1)
	}

	rl.Shutdown()
}
