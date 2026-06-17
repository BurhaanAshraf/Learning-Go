package main

import (
	"context"
	"fmt"
	"time"
)

// Token Bucket rate limiter.
//
// Concept: the bucket holds up to N tokens. Each request consumes one
// token; an empty bucket means the request is rejected. Tokens refill at a
// fixed rate over time. This allows short bursts (up to N requests at
// once) while still capping the long-term average rate to the refill rate.
//
// Go mapping:
//   bucket       -> buffered channel, make(chan struct{}, N)
//   token        -> one value in the channel (struct{} is 0 bytes — we
//                    only care that a slot is occupied, not its value)
//   take a token -> receive from the channel
//   bucket empty -> channel has nothing to receive
//   refill       -> background goroutine sends into the channel on a timer
//   bucket full  -> the send is dropped (non-blocking) instead of blocking
//
// Two independent, often-confused knobs:
//   burst capacity = bucket size (channel capacity)      -> how spiky traffic can be
//   sustained rate = refill rate (channel send interval) -> long-term average throughput
// Example: size=10, refill=1/sec -> 10 requests succeed instantly, then
// throughput is capped to 1/sec until the bucket refills.
//
// Refill strategy: this implementation uses a ticker (add one token per
// tick). Simple, but a tick missed during a goroutine stall (GC pause,
// scheduler delay) is a token lost permanently. Production limiters
// instead compute tokens from elapsed wall time
// (tokensToAdd = elapsed.Seconds() * refillRatePerSecond), which
// self-corrects after a stall instead of losing capacity. Fine to use a
// ticker for low-stakes/internal use; know the tradeoff before relying on
// this for anything where lost tokens matter.
//
// Performance Caveat: Avoid setting refillTime to ultra-low intervals (e.g., microseconds).
// The system overhead of the ticker and context switching will cause the channel
// writes to lag behind real-time, artificially throttling throughput.
//
// Concurrency: channels are safe for concurrent send/receive, so Allow()
// can be called from many goroutines with no extra locking.
//
// Related algorithms, for context: Leaky Bucket processes at a strictly
// fixed rate with no burst allowance. Fixed Window Counter resets a
// counter every interval (simple, but allows double-bursts at window
// edges). Sliding Window avoids that edge case by moving the window
// continuously instead of resetting it.

type RateLimiter struct {
	// Thread-safe token store. Kept internal to prevent external manual refills.
	tokens     chan struct{}
	refillTime time.Duration
	cancel     context.CancelFunc
}

// newRateLimiter creates a limiter allowing `events` requests in a burst,
// refilling one token every refillTime. Starts full so the first burst
// succeeds immediately.
func newRateLimiter(events int, refillTime time.Duration) *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())

	rl := &RateLimiter{
		tokens:     make(chan struct{}, events),
		refillTime: refillTime,
		cancel:     cancel,
	}

	// Seed the bucket to full capacity initially.
	for range events {
		rl.tokens <- struct{}{}
	}

	go rl.refill(ctx)
	return rl
}

// refill adds one token per tick. The inner select+default is required:
// without it, sending to an already-full channel blocks forever, which
// would permanently stall this goroutine even after tokens are later
// consumed and room opens up again.
func (rl *RateLimiter) refill(ctx context.Context) {
	ticker := time.NewTicker(rl.refillTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Channel is full; drop the token. Prevents the background
				// goroutine from blocking indefinitely.
			}
		}
	}
}

// Stop cancels the refill goroutine. Always call this when a RateLimiter
// is no longer needed. Without it, the goroutine leaks for the program's
// lifetime — harmless for one long-lived limiter in main(), but a real
// leak if limiters are created and discarded per request.
func (rl *RateLimiter) Stop() {
	if rl.cancel != nil {
		rl.cancel()
	}
}

// Allow checks for a token without blocking.
// true  -> token taken, request proceeds.
// false -> bucket empty, reject now.
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Wait blocks until a token becomes available or the context is canceled.
// Use this variant for background workers or queue consumers where delaying
// execution is preferable to outright rejection.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func TokenBucket_RL() {
	rl := newRateLimiter(10, 1*time.Second) // burst 10, refill 1/sec
	defer rl.Stop()

	for i := 0; i < 20; i++ {
		if rl.Allow() {
			fmt.Println("Token Accepted")
		} else {
			fmt.Println("Token Rejected")
		}
		time.Sleep(300 * time.Millisecond)
	}
}
