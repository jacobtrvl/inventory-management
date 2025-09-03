// Copyright 2025 Jacob Philip. All rights reserved.
// Super simple rate limiter using a token bucket algorithm.
// Requests to Allow waits until a token is available.
// No busy looping! Uses the magic of Go channels!
package ratelimiter

import (
	"time"
)

const (
	maxQueueSize = 10000
	maxOpsPerSec = 1000000
)

type RateLimiter struct {
	operationsPerSecond int
	tokens              int
	check               chan chan struct{}
	done                chan struct{}
	requests            []chan struct{}
}

func NewRateLimiter(o int) *RateLimiter {
	if o <= 0 {
		o = 1
	}
	if o > maxOpsPerSec {
		o = maxOpsPerSec
	}
	rl := &RateLimiter{
		operationsPerSecond: o,
		tokens:              o,
		check:               make(chan chan struct{}),
		done:                make(chan struct{}),
		requests:            make([]chan struct{}, 0),
	}
	go rl.run()
	return rl
}

func (rl *RateLimiter) run() {
	tokenAllot := time.NewTicker(time.Second / time.Duration(rl.operationsPerSecond))
	heartBeat := time.NewTicker(time.Millisecond)
	defer tokenAllot.Stop()
	defer heartBeat.Stop()

	for {
		select {
		case <-heartBeat.C:
			for rl.tokens > 0 && len(rl.requests) > 0 {
				req := rl.requests[0]
				rl.requests[0] = nil
				rl.requests = rl.requests[1:]
				req <- struct{}{}
				rl.tokens--
			}
		case <-tokenAllot.C:
			if rl.tokens < rl.operationsPerSecond {
				rl.tokens++
			}
		case <-rl.done:
			return
		case req := <-rl.check:
			if len(rl.requests) > maxQueueSize {
				continue // Drop request if too many pending
			}
			rl.requests = append(rl.requests, req)

		}
	}
}

func (rl *RateLimiter) Allow() {
	req := make(chan struct{})
	rl.check <- req
	<-req
}

func (rl *RateLimiter) Shutdown() {
	close(rl.done)
}
