package main

import (
	"fmt"
	"sync"
	"time"
)

// Ratelimiter implements a thread-safe Fixed Window Counter rate limiter.
//
// CONCURRENCY MECHANICS:
// Unlike a Token Bucket that can store capacity as signaling values in a buffered
// channel to achieve free thread safety, a Fixed Window tracking state consists of
// plain primitive types (int and time.Time). Because the pattern "check count, then
// increment" is non-atomic, a sync.Mutex is strictly required to prevent race conditions
// where concurrent goroutines read stale state simultaneously.
//
// IDLE BEHAVIOR:
// There is no retention of unused capacity. Being idle for multiple windows does not
// bank tokens for later use; the user is strictly capped at the configured limit
// within any active window.
//
// ARCHITECTURAL SIMPLICITY:
// Because this algorithm uses lazy evaluation on incoming requests to advance windows,
// it eliminates the need for a background ticker goroutine. This completely removes
// any risk of resource or memory leaks, meaning it does not require a Stop() cleanup method.
type Ratelimiter struct {
	mu        sync.Mutex
	count     int
	limit     int
	window    time.Duration
	resetTime time.Time // Uses zero-value initialization trick (time.Time{} / Year 1)
}

// NewRateLimiter initializes a rate limiter with a maximum request limit
// allowed per window duration.
func NewRateLimiter(limit int, window time.Duration) *Ratelimiter {
	return &Ratelimiter{
		limit:  limit,
		window: window,
	}
}

// Allow atomically determines if an incoming request falls within the current quota.
//
// ZERO-VALUE WINDOW INITIALIZATION:
// On the very first request, now.After(rl.resetTime) evaluates to true instantly
// because rl.resetTime defaults to the zero-value timestamp. This lazily sets the
// first true deadline and zeroes the counter without requiring explicit setup code.
//
// KNOWN WEAKNESS (THE BOUNDARY BURST):
// This implementation is vulnerable to a 2x burst capability right at window edges.
// If a user exhausts their entire limit at the tail-end of Window N, and then exhausts
// it again at the absolute beginning of Window N+1, they can successfully push
// (2 * limit) requests within a sub-window timeframe. Use a Sliding Window Counter
// if strict traffic smoothing is required.
func (rl *Ratelimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// If the current window boundary has passed, cycle to the next window
	if now.After(rl.resetTime) {
		rl.resetTime = now.Add(rl.window)
		rl.count = 0
	}

	// Evaluate against the boundary capacity
	if rl.count < rl.limit {
		rl.count++
		return true
	}
	return false
}

func fixedWindow_RL() {
	// Enforce a strict quota of 5 requests per 2 seconds
	rateLimiter := NewRateLimiter(5, 2*time.Second)

	var wg sync.WaitGroup

	// Simulate a sudden concurrent spike of 10 incoming requests.
	// Out of these 10 concurrent routines, exactly 5 will be allowed and 5 rejected.
	// The precise order of success is non-deterministic due to Go runtime scheduling.
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rateLimiter.Allow() {
				fmt.Println("Token Accepted")
			} else {
				fmt.Println("Token Rejected")
			}
		}()
	}
	wg.Wait()
}
