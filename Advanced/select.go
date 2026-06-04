package main

import (
	"fmt"
	"time"
)

/*
=============================================================================
SELECT STATEMENT & ADVANCED CHANNEL PATTERNS — COMPLETE REVISION NOTES
=============================================================================

─────────────────────────────────────────────────────────────────────────────
WHAT IS SELECT?
─────────────────────────────────────────────────────────────────────────────

select is Go's way of handling multiple channels at once. It is like a
switch statement, but for channel operations. It waits until one of its
cases is ready, then executes that case.

    select {
    case msg := <-ch1:
        runs if ch1 has a value
    case msg := <-ch2:
        runs if ch2 has a value
    }

If multiple cases are ready at the same time, Go picks one at random.
This is intentional — no channel gets priority over another.

─────────────────────────────────────────────────────────────────────────────
BLOCKING SELECT vs NON-BLOCKING SELECT
─────────────────────────────────────────────────────────────────────────────

WITHOUT default:
    select blocks until at least one case is ready.
    It will wait as long as needed — could be milliseconds or forever.
    Use this when you WANT to wait for a channel.

    select {
    case msg := <-ch1:
        fmt.Println(msg)
    case msg := <-ch2:
        fmt.Println(msg)
    }
    blocks here until ch1 or ch2 has a value

WITH default:
    select checks all cases instantly. If none are ready, it runs default
    immediately and moves on — does NOT wait for any channel.
    Use this when you do NOT want to block — you want to check
    opportunistically and continue if nothing is ready.

    select {
    case msg := <-ch1:
        fmt.Println(msg)
    default:
        fmt.Println("nothing ready, moving on")
    }
    never blocks — runs default if ch1 is empty

RULE:
    No default  = blocking select  → waits for a channel to be ready
    With default = non-blocking select → checks and moves on immediately

─────────────────────────────────────────────────────────────────────────────
TOPIC 1: NON-BLOCKING SELECT WITH OPPORTUNISTIC MATCHING
─────────────────────────────────────────────────────────────────────────────

The pattern: launch goroutines to fill buffered channels, give them a moment
to run, then use select with default to opportunistically drain them.

Buffered channels are used here so goroutines can send without waiting for
a receiver to be ready. They just drop the value in the buffer and signal
completion via channelSync.

The time.Sleep(10ms) is used to give the Go scheduler time to run the
goroutines before the select loop executes. Without it, select would run
instantly before goroutines get scheduled, hit default on every iteration,
and never collect the values.

IMPORTANT NUANCE:
    The 10ms sleep is a very safe assumption, NOT a hard guarantee.
    On an extremely loaded system, goroutines could theoretically not
    finish in 10ms. For production code, use sync.WaitGroup or a
    proper done channel instead of relying on sleep timing.

HOW THE LOOP WORKS:
    We loop exactly 2 times because we know there are exactly 2 values.
    Each iteration, select picks whichever channel has a value ready.
    If both are ready, Go picks one at random — order is not guaranteed.
    After both values are collected, we drain channelSync to confirm
    both goroutines have fully completed their lifecycle.

─────────────────────────────────────────────────────────────────────────────
TOPIC 2: TIMEOUTS USING time.After
─────────────────────────────────────────────────────────────────────────────

time.After(d) returns a read-only channel (<-chan time.Time) that the Go
runtime automatically sends a timestamp to after duration d expires.
You never send to it — the runtime does. You just receive from it.

This makes it perfect as a select case — if your worker channel does not
produce a value before the timeout fires, the timeout case wins.

    select {
    case msg := <-ch:
        worker finished in time
    case <-time.After(2 * time.Second):
        timeout fired first — worker took too long
    }

GOROUTINE LEAK WARNING:
    When the timeout case wins, the worker goroutine is still running in
    the background. It will eventually try to send to ch but nobody is
    receiving anymore. That goroutine is now leaked — it is stuck forever.

    In production code, use a context with cancellation to signal the
    worker to stop when timeout occurs:

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    pass ctx to worker so it can stop when cancelled

    For learning purposes, the time.After pattern is fine to understand
    the concept. Just be aware of the leak in real code.

─────────────────────────────────────────────────────────────────────────────
TOPIC 3: CLOSED CHANNEL LIFECYCLE & COMMA-OK IDIOM
─────────────────────────────────────────────────────────────────────────────

When a channel is closed:
    → Remaining values in the buffer can still be received normally
    → Once buffer is empty, any further receive returns instantly with
      the zero value of the type (0 for int, "" for string, false for bool)
    → The channel does NOT block after being closed — it returns immediately

The comma-ok idiom lets you detect whether a channel is closed:

    msg, ok := <-ch

    ok == true  → channel is open, msg is a real value
    ok == false → channel is closed and empty, msg is zero value

HOW THE ITERATION WORKS (depends on how many values are in channel):

    If channel has 1 value (ch <- 20) then close(ch):
        Iteration 1: receives 20, ok = true  → prints value
        Iteration 2: channel closed and empty, ok = false → cleanup and return

    If channel had 3 values then close(ch):
        Iterations 1,2,3: receive real values, ok = true
        Iteration 4: ok = false → cleanup and return

    The number of ok=true iterations always equals the number of values
    sent before close. After that, one final ok=false iteration triggers cleanup.

WHY return INSIDE select INSIDE for:
    return exits the entire function, which tears down the infinite for loop.
    This is cleaner than using a break + boolean flag.
    break alone inside a select only breaks out of the select, not the for loop.

─────────────────────────────────────────────────────────────────────────────
STRUCT{} AS A SIGNAL — WHY NOT bool OR int?
─────────────────────────────────────────────────────────────────────────────

    channelSync := make(chan struct{})
    channelSync <- struct{}{}

struct{} is an empty struct — it has zero size and allocates no memory.
When you only need a signal (I am done) and the value does not matter,
struct{} is the most efficient and idiomatic choice in Go.

    bool uses 1 byte  — value does not matter but memory is wasted
    int  uses 8 bytes — value does not matter but memory is wasted
    struct{} uses 0 bytes — perfect for pure signaling

This is a common Go convention you will see in real codebases.

─────────────────────────────────────────────────────────────────────────────
SUMMARY — KEY RULES TO REMEMBER
─────────────────────────────────────────────────────────────────────────────

1. select handles multiple channels — runs whichever case is ready first.
2. If multiple cases ready simultaneously — Go picks one at RANDOM.
3. select without default — BLOCKS until a case is ready.
4. select with default — NON-BLOCKING, runs default if nothing is ready.
5. time.After returns a channel the runtime sends to after a duration.
6. time.After timeout can cause goroutine leaks — use context in production.
7. Closed channel does not block — returns zero value + ok=false instantly.
8. comma-ok idiom (msg, ok := <-ch) detects if channel is closed.
9. break inside select only breaks select, not the enclosing for loop.
10. struct{} is the idiomatic signal type — zero size, zero allocation.

=============================================================================
*/

func Select() {
	fmt.Println("=== Topic 1: Non-Blocking Select with Opportunistic Matching ===")
	nonBlockingSelect()

	fmt.Println("\n=== Topic 2: Timeout with time.After ===")
	timeAfter()

	fmt.Println("\n=== Topic 3: Comma-Ok Idiom & Closed Channel Lifecycle ===")
	ok()
}

// =============================================================================
// TOPIC 1 — NON-BLOCKING SELECT WITH OPPORTUNISTIC MATCHING
// =============================================================================
//
// Pattern: fill buffered channels via goroutines, give scheduler time to run
// them, then drain channels using select with default (non-blocking).
//
// Buffered channels (capacity 1) are used so goroutines can send immediately
// without waiting for a receiver — they drop value in buffer and move on.
//
// channelSync (capacity 2) tracks full goroutine lifecycle completion.
// We receive from it twice at the end to confirm both goroutines are truly done.

func nonBlockingSelect() {

	ch1 := make(chan int, 1)     // buffered — goroutine sends without blocking
	ch2 := make(chan int, 1)     // buffered — goroutine sends without blocking
	channelSync := make(chan struct{}, 2) // buffered — tracks goroutine completion

	go func() {
		ch1 <- 10              // drops into buffer, does not block
		channelSync <- struct{}{} // signals this goroutine is fully done
	}()

	go func() {
		ch2 <- 20              // drops into buffer, does not block
		channelSync <- struct{}{} // signals this goroutine is fully done
	}()

	// Give the Go scheduler time to run both goroutines before select executes.
	// Without this, select runs instantly before goroutines are scheduled,
	// hits default on every iteration, and never collects the values.
	// NOTE: this is a safe assumption for learning, not a hard guarantee.
	// In production, use sync.WaitGroup or a proper done channel instead.
	time.Sleep(10 * time.Millisecond)

	// Loop exactly 2 times — one per value we know exists in the channels.
	// Each iteration, select picks whichever channel has a value ready.
	// If both are ready simultaneously, Go picks one at RANDOM — no priority.
	// default case would run if neither channel had a value (only if sleep removed).
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println("Received Message From Channel 1:", msg)
		case msg := <-ch2:
			fmt.Println("Received Message from Channel 2:", msg)
		default:
			// only reachable if sleep is removed and goroutines not yet scheduled.
			// with the 10ms sleep above, both channels are already populated
			// so this default case will NOT run.
			fmt.Println("No Message Received...")
		}
	}

	// drain channelSync twice — confirms both goroutines have fully completed.
	// without this, main could exit before goroutines finish their cleanup work.
	<-channelSync
	<-channelSync
}

// =============================================================================
// TOPIC 2 — OPPORTUNISTIC TIMEOUTS USING time.After
// =============================================================================
//
// Pattern: race a worker channel against a timeout channel using select.
// Whichever is ready first wins. The other case is abandoned.
//
// time.After(d) returns a <-chan time.Time that the Go runtime automatically
// sends a timestamp to after duration d. You never send to it — runtime does.
//
// GOROUTINE LEAK: when timeout wins, the worker goroutine is still sleeping
// and will eventually try to send to ch with nobody receiving. It leaks.
// For production code, use context.WithTimeout to cancel the worker properly.

func timeAfter() {

	ch := make(chan int) // unbuffered — worker and receiver must meet

	go func() {
		// worker takes 3 seconds — longer than the 2 second timeout below.
		// this guarantees the timeout case wins in select.
		time.Sleep(3 * time.Second)
		ch <- 10 // this send will never complete — nobody receiving after timeout
		         // this goroutine is now leaked (stuck here forever)
	}()

	// select races ch against the timeout channel.
	// worker needs 3s but timeout fires at 2s → timeout case wins.
	// if worker was faster (e.g. 1s sleep), the msg case would win instead.
	select {
	case msg := <-ch:
		// only reaches here if worker sends before 2 second timeout
		fmt.Println("Value Received:", msg)
	case <-time.After(2 * time.Second):
		// reaches here because worker took 3s but timeout fired at 2s
		fmt.Println("Timeout! Worker took too long.")
	}
	// NOTE: goroutine is still running here, leaked, stuck on ch <- 10 forever.
	// in production, signal cancellation via context before returning.
}

// =============================================================================
// TOPIC 3 — CLOSED CHANNEL LIFECYCLE & COMMA-OK IDIOM
// =============================================================================
//
// Pattern: receive from a channel in an infinite loop using select + comma-ok.
// Stop cleanly when the channel is closed and drained.
//
// What happens when channel is closed:
//   → Values already in buffer can still be received normally (ok = true)
//   → Once buffer is empty, receive returns instantly with zero value (ok = false)
//   → Channel does NOT block after being closed — it unblocks immediately
//
// Comma-ok idiom:
//   msg, ok := <-ch
//   ok = true  → channel open, msg is real value
//   ok = false → channel closed and empty, msg is zero value (0 for int)
//
// Why return instead of break:
//   break inside select only exits the select block, NOT the for loop.
//   return exits the entire function, cleanly terminating the infinite loop.

func ok() {

	ch := make(chan int) // unbuffered — goroutine and receiver meet directly

	go func() {
		ch <- 20   // sends one value — receiver gets this with ok = true
		close(ch)  // closes channel — next receive will get ok = false
		           // closing signals "no more values will ever be sent"
	}()

	// infinite loop — keeps receiving until channel is closed and empty.
	// select here allows extending this pattern to multiple channels easily.
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				// channel is closed and empty.
				// msg is 0 (zero value for int) — do not use it.
				// ok = false is the signal to stop and clean up.
				fmt.Println("Channel Closed! Executing structural cleanup...")
				return // exits the function, terminating the infinite for loop
				       // break here would only exit select, not the for loop
			}
			// ok = true — msg is a real value sent by the goroutine
			fmt.Println("Values Received:", msg)
		}
	}
}