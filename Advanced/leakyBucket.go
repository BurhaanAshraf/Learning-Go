package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// LEAKY BUCKET ALGORITHM — REVISION NOTES
// ============================================================================
//
// WHAT
//   A rate limiting algorithm that allows capacity to recover at a fixed rate
//   over time. Requests consume capacity; capacity recovers steadily.
//
// WHY
//   Same goal as Token Bucket — limit how many requests a system handles.
//   Different in HOW capacity is tracked and recovered.
//
// ============================================================================
// LEAKY BUCKET VS TOKEN BUCKET
// ============================================================================

// ----------------------------------------------------------------------------
// 1. TOKEN BUCKET (The Accumulation Model)
// ----------------------------------------------------------------------------
// - Core Metric: Tracks AVAILABLE TOKENS. Time adds tokens; requests subtract them.
// - Traffic Shape: ALLOWS BURSTS. If the bucket is full, a massive wave of concurrent
//   requests passes through instantly until the capacity is completely drained.
//
// ----------------------------------------------------------------------------
// 2. LEAKY BUCKET (The Traffic-Shaping Queue)
// ----------------------------------------------------------------------------
// - Core Metric: Tracks USED CAPACITY (Water/Space). Time drains water; requests add it.
// - Traffic Shape: SMOOTHS BURSTS. Even if a wave of requests hits simultaneously,
//   they are serialized into a steady, constant, predictable output drip.
// ============================================================================
//
//   Core distinction:
//     Token Bucket → goroutine pushes tokens in on a schedule.
//     Leaky Bucket → tokens are calculated from elapsed time on demand.
//
//   Both allow bursting up to capacity. The practical differences are
//   in implementation complexity, goroutine lifecycle, and persistence
//   (covered in Real-World Use below).
//
// ============================================================================
// HOW TIME-BASED REFILL WORKS (THE CORE MECHANIC)
// ============================================================================
//
// Every time Allow() is called, it calculates how many tokens should have
// recovered since the last call, adds them, then checks if one is available.
//
// STEP BY STEP inside Allow():
//
//   1. elapsed = now - lastLeakTime
//
//   2. tokensToAdd = elapsed / leakRate   (integer division)
//      Only whole intervals count. 1.3 intervals → 1 token added.
//
//   3. tokens = min(tokens + tokensToAdd, capacity)
//      Never exceed capacity — idle time can accumulate at most `capacity` tokens.
//
//   4. lastLeakTime += tokensToAdd * leakRate
//
//      WHY NOT lastLeakTime = now?
//        elapsed = 1.3 * leakRate → tokensToAdd = 1 (0.3 dropped by int division).
//        Setting lastLeakTime = now discards that 0.3 permanently.
//        Advancing by (tokensToAdd * leakRate) preserves the 0.3 remainder,
//        so it counts toward the next call. This keeps token accounting accurate
//        over many calls.
//
//   5. If tokens > 0: consume one → return true.
//      If tokens == 0: return false.
//
// ============================================================================
// WHY A MUTEX IS NEEDED HERE
// ============================================================================
//
// Token Bucket uses a channel — goroutine-safe by design, no mutex needed.
//
// Leaky Bucket uses plain integer fields (tokens, lastLeakTime). Without a
// mutex, two goroutines could both read tokens = 1, both pass the check,
// both decrement — two requests pass when only one should have.
//
// The mutex makes Allow() fully atomic: one goroutine at a time reads state,
// computes the refill, checks, and decrements.
//
// ============================================================================
// COMMON PITFALLS
// ============================================================================
//
// 1. Setting lastLeakTime = now instead of advancing it by (tokensToAdd * leakRate)
//    Discards the fractional interval on every call — tokens are issued
//    slightly less often than they should be over time.
//
// 2. Not capping tokens at capacity
//    After a long idle period, tokensToAdd can exceed capacity.
//    Without the cap, tokens grow beyond the intended limit.
//    Always: if tokens > capacity { tokens = capacity }
//
// 3. Reading tokens or lastLeakTime outside the mutex
//    Both fields are modified together in Allow(). Reading either one without
//    holding the lock is a data race — even for a read-only check.
//
// ============================================================================
// REAL-WORLD USE
// ============================================================================
//
//   Leaky Bucket is often preferred over Token Bucket in production because:
//
//   - No background goroutine — simpler lifecycle, nothing to leak or cancel.
//
//   - Handles scheduler delays automatically — if the program is paused for
//     3 seconds, the next Allow() call computes 3s worth of tokens at once
//     from elapsed time. A ticker-based approach may miss those ticks.
//
//   - Distributable — the state is just two values: tokens (int) and
//     lastLeakTime (time.Time). These can be stored in Redis or a DB.
//     On each request: load state, compute refill, check, decrement, save.
//     A channel-based Token Bucket cannot be distributed this way.
//
//   Common uses:
//     - Per-user API rate limiting in HTTP middleware.
//     - Database query throttling.
//     - Outgoing webhook delivery rate control.
//
// ============================================================================
// WHEN TO USE WHICH
// ============================================================================
//
//   Token Bucket (channel-based):
//     - Simple, idiomatic Go.
//     - Single-process only.
//     - Fine with a background goroutine running.
//
//   Leaky Bucket (time-based):
//     - No goroutine needed.
//     - State can be persisted for distributed rate limiting.
//     - More accurate under real-world scheduling pressure.
//
// ============================================================================
// INTERVIEW POINTS
// ============================================================================
//
//   Q: What is the difference between Token Bucket and Leaky Bucket?
//   A: Token Bucket uses a goroutine to add tokens periodically. Leaky
//      Bucket computes tokens from elapsed time on each Allow() call.
//      Both allow bursting; Leaky Bucket needs no goroutine and can be
//      distributed by persisting its two-field state.
//
//   Q: Why advance lastLeakTime by (tokensToAdd * leakRate) instead of now?
//   A: Integer division drops fractional intervals. Advancing by the exact
//      amount covered preserves the remainder for the next call. Setting
//      it to now discards the remainder permanently, causing under-refilling.
//
//   Q: Why does Leaky Bucket need a mutex but Token Bucket doesn't?
//   A: Token Bucket uses a channel — goroutine-safe by design. Leaky Bucket
//      uses plain integer fields — concurrent reads and writes on those
//      without a mutex are a data race.
//
//   Q: Why is Leaky Bucket preferred for distributed systems?
//   A: State is just tokens (int) + lastLeakTime (time.Time) — both can be
//      stored in Redis or a DB and computed identically on any server.
//      A channel cannot be shared across processes.
//
// ============================================================================

// leakyBucket holds the state for one rate limiter instance.
//
//	capacity     — max tokens the bucket can hold.
//	leakRate     — time required for one token to recover.
//	               e.g. 500ms → 2 tokens recover per second.
//	tokens       — tokens currently available.
//	lastLeakTime — timestamp of the last token calculation.
//	               Used to compute elapsed time on the next Allow() call.
//	mu           — protects tokens and lastLeakTime from concurrent access.
type leakyBucket struct {
	capacity     int
	leakRate     time.Duration
	tokens       int
	lastLeakTime time.Time
	mu           sync.Mutex
}

// NewLeakyBucket creates a bucket that starts full.
// All capacity is immediately available for the first burst.
func NewLeakyBucket(capacity int, leakRate time.Duration) *leakyBucket {
	return &leakyBucket{
		capacity:     capacity,
		leakRate:     leakRate,
		tokens:       capacity,
		lastLeakTime: time.Now(),
	}
}

// Allow returns true if a token is available (request proceeds),
// false if the bucket is empty (request rejected).
// Safe to call from multiple goroutines concurrently.
func (lb *leakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastLeakTime)

	// How many full leakRate intervals have passed?
	// Integer division drops the fractional remainder — it carries forward.
	tokensToAdd := int(elapsed / lb.leakRate)

	if tokensToAdd > 0 {
		lb.tokens += tokensToAdd
		if lb.tokens > lb.capacity {
			lb.tokens = lb.capacity
		}
		// WHY WE ADVANCE lastLeakTime BY (tokensToAdd * leakRate), NOT by setting it to now
		//
		// lastLeakTime is an anchor — it marks "where we last counted from."
		// Every time Allow() runs, we measure how much time passed since that anchor,
		// and move the anchor forward by exactly the time we used up.
		//
		// The problem with lastLeakTime = now:
		//   Suppose leakRate = 1s and a request arrives at t=1.3s.
		//   elapsed = 1.3s → tokensToAdd = 1 (integer division drops the 0.3s)
		//   If we set lastLeakTime = now (1.3s), the 0.3s is thrown away forever.
		//
		// The next request arrives at t=2.0s:
		//   elapsed = 2.0 - 1.3 = 0.7s → tokensToAdd = 0 (not enough for a full token)
		//   Request rejected, even though 2.0s of real time has passed since t=0!
		//
		// The fix — advance by (tokensToAdd * leakRate):
		//   tokensToAdd = 1, leakRate = 1s → advance anchor by 1.0s → anchor is now 1.0s
		//   The 0.3s remainder is NOT discarded — it's still "ahead of the anchor."
		//
		// The next request arrives at t=2.0s:
		//   elapsed = 2.0 - 1.0 = 1.0s → tokensToAdd = 1 → request accepted ✓
		//
		// TIMELINE (leakRate = 1s):
		//
		//   t=0.0s         t=1.0s    t=1.3s    t=2.0s
		//    |              |         |          |
		//    anchor         |      req arrives   next req
		//                   |
		//              correct anchor after req 1
		//                   (not 1.3s — that wastes the 0.3s)
		//
		//   correct:  anchor moves to 1.0s → next elapsed = 1.0s → 1 token ✓
		//   wrong:    anchor moves to 1.3s → next elapsed = 0.7s → 0 tokens ✗

		lb.lastLeakTime = lb.lastLeakTime.Add(time.Duration(tokensToAdd) * lb.leakRate)
	}

	if lb.tokens > 0 {
		lb.tokens--
		return true
	}
	return false
}

func leakyBucket_RL() {
	// capacity 5: up to 5 requests allowed in a burst.
	// leakRate 500ms: 1 token recovers every 500ms (2 per second).
	lb := NewLeakyBucket(5, 500*time.Millisecond)

	var wg sync.WaitGroup

	// Burst 1: 10 requests hit at once. Bucket starts full at 5.
	// Expected: 5 accepted, 5 rejected.
	fmt.Println("=== BURST 1: bucket full (5 tokens) ===")
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if lb.Allow() {
				fmt.Println("Accepted")
			} else {
				fmt.Println("Rejected")
			}
		}()
	}
	wg.Wait()

	// Wait 1 second — at 1 token per 500ms, 2 tokens recover.
	fmt.Println("\nWaiting 1s — expecting 2 tokens to recover...")
	time.Sleep(1 * time.Second)

	// Burst 2: 10 requests. Bucket should have 2 tokens.
	// Expected: 2 accepted, 8 rejected.
	fmt.Println("\n=== BURST 2: 2 tokens recovered ===")
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if lb.Allow() {
				fmt.Println("Accepted")
			} else {
				fmt.Println("Rejected")
			}
		}()
	}
	wg.Wait()

	fmt.Println("\nDone.")
}
