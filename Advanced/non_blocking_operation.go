package main

import (
	"fmt"
	"time"
)

/*
=============================================================================
NON-BLOCKING CHANNEL OPERATIONS WITH SELECT — REVISION NOTES
=============================================================================

─────────────────────────────────────────────────────────────────────────────
HOW select EVALUATES — THE CORE MECHANIC
─────────────────────────────────────────────────────────────────────────────

select checks all its cases simultaneously in a single instant:

    1. ONE case ready      → execute that case immediately
    2. MULTIPLE ready      → Go picks one at RANDOM (no priority)
    3. NONE ready          → execute default immediately (if present)
    4. NONE ready, no default → block until a case becomes ready

select itself never sleeps or waits. It evaluates and picks in one shot.
Any sleeping or waiting happens INSIDE the case/default blocks — not in
select. This distinction matters for understanding the sleep trade-off
in Pattern 3.

─────────────────────────────────────────────────────────────────────────────
select WITH default vs WITHOUT default
─────────────────────────────────────────────────────────────────────────────

WITH default    → NON-BLOCKING. If no case is ready, default runs instantly.
WITHOUT default → BLOCKING. select waits until at least one case is ready.

─────────────────────────────────────────────────────────────────────────────
PATTERN 1 — NON-BLOCKING RECEIVE
─────────────────────────────────────────────────────────────────────────────

    ch := make(chan int, 1)
    ch <- 5  // value already in buffer

    select {
    case msg := <-ch:   // ch has a value → case is ready → runs instantly
    default:            // skipped — a ready case was found
    }

    If ch were EMPTY:
    select {
    case msg := <-ch:   // ch is empty → case NOT ready → skipped
    default:            // runs instantly
    }

─────────────────────────────────────────────────────────────────────────────
PATTERN 2 — NON-BLOCKING SEND
─────────────────────────────────────────────────────────────────────────────

    ch := make(chan int)  // unbuffered

    select {
    case ch <- 1:   // only ready if another goroutine is actively waiting
                       to receive right now. nobody waiting → not ready.
    default:        // runs instantly — data dropped, no deadlock, no blocking
    }

Without default, ch <- 1 would block forever since nobody is receiving.
With default, it just skips gracefully. This is how you drop data safely
when the receiver is not ready — no panic, no freeze.

─────────────────────────────────────────────────────────────────────────────
PATTERN 3 — WORKER LOOP WITH GRACEFUL SHUTDOWN
─────────────────────────────────────────────────────────────────────────────

Three channels, three distinct roles:

    data  →  incoming values from main to worker
    quit  →  main signals worker to stop
    done  →  worker confirms to main that it fully stopped

WHY quit AND done — not just quit alone?

    quit <- true   sends the intent to stop.
    <-done         confirms the stop actually happened.

    Without <-done: main exits right after quit <- true. Go runtime kills
    all goroutines instantly — worker is aborted mid-cleanup, never finishes.
    <-done forces main to block until worker explicitly says "I am done."

    quit = "please stop"
    done = "I have stopped"

WHY sleep IN default — AND THE TRADE-OFF:

    select picks default instantly when data and quit are both empty.
    THEN the default block runs its sleep. The sleep is not inside select —
    select already finished. This matters because:

    During the sleep, if data arrives on the channel, it CANNOT be processed
    until the sleep ends and the loop cycles back to select again.

    PRO: prevents the for loop from spinning at 100% CPU when idle
    CON: introduces up to 100ms latency before new data is processed

    WRONG: "select sleeps for 100ms"
    RIGHT:  "select picks default instantly, THEN default sleeps 100ms"

─────────────────────────────────────────────────────────────────────────────
SUMMARY — KEY RULES
─────────────────────────────────────────────────────────────────────────────

1. select with default    = non-blocking (skips if not ready)
2. select without default = blocking (waits until ready)
3. Multiple cases ready   = Go picks at RANDOM
4. sleep in default does NOT pause select — select already finished
5. sleep in default = CPU savings but adds latency up to sleep duration
6. quit signals intent to stop. done confirms stop happened.
7. Without <-done, main exits and kills worker mid-cleanup

=============================================================================
*/

func NonBlocking() {

	// =========================================================================
	// PATTERN 1 — NON-BLOCKING RECEIVE
	// =========================================================================
	// ch1 already has a value in its buffer before select runs.
	// receive case is immediately ready → executes, default is skipped.
	// if ch1 were empty, receive case would not be ready → default would run.
	// =========================================================================

	ch1 := make(chan int, 1) // buffered, capacity 1
	ch1 <- 5                 // value in buffer — no receiver needed

	select {
	case msg := <-ch1:
		// ch1 has data → ready → runs instantly. default skipped.
		fmt.Println("Received:", msg)
	default:
		// only runs if ch1 is empty. skipped here.
		fmt.Println("No message available!")
	}

	// =========================================================================
	// PATTERN 2 — NON-BLOCKING SEND
	// =========================================================================
	// Unbuffered channel, no goroutine waiting to receive.
	// send case is not ready → default runs instantly.
	// without default, ch <- 1 would block forever → deadlock.
	// with default, data is dropped gracefully — no blocking, no panic.
	// =========================================================================

	ch := make(chan int) // unbuffered — sender and receiver must meet

	select {
	case ch <- 1:
		// only runs if a goroutine is actively waiting to receive right now.
		// nobody is waiting → not ready → skipped.
		fmt.Println("Value Sent")
	default:
		// no receiver ready → send case skipped → default runs instantly.
		// data is dropped gracefully instead of causing a deadlock.
		fmt.Println("No receiver ready — value dropped gracefully")
	}

	// =========================================================================
	// PATTERN 3 — WORKER LOOP WITH GRACEFUL SHUTDOWN
	// =========================================================================
	// Worker runs an infinite for-select loop:
	//   data case    → process incoming value immediately when ready
	//   quit case    → begin shutdown, signal done, return
	//   default case → neither data nor quit ready, sleep and loop again
	//
	// data → main sends values to worker
	// quit → main tells worker to stop
	// done → worker tells main it has fully stopped (graceful shutdown)
	// =========================================================================

	data := make(chan int)  // unbuffered — main and worker meet on each send
	quit := make(chan bool) // unbuffered — shutdown signal
	done := make(chan bool) // unbuffered — worker confirms full stop to main

	go func() {
		for {
			select {
			case d := <-data:
				// data is ready → process immediately.
				// select picks this the instant main sends a value.
				fmt.Println("Data Received:", d)

			case <-quit:
				// shutdown signal received → clean up and confirm to main.
				fmt.Println("Stopping Worker safely...")
				done <- true // tell main: fully stopped, safe to exit now
				return       // exit goroutine

			default:
				// data and quit both empty → select picks default instantly.
				// sleep here prevents 100% CPU spinning when worker is idle.
				// TRADE-OFF: new data arriving during sleep waits up to 100ms
				// before being processed — select only re-evaluates after sleep ends.
				time.Sleep(100 * time.Millisecond)
				fmt.Println("Waiting For Data...")
			}
		}
	}()

	// give worker time to start and execute default at least once
	time.Sleep(150 * time.Millisecond)

	// send 5 values — each blocks until worker receives (unbuffered rendezvous)
	for i := 0; i < 5; i++ {
		data <- i
	}

	// tell worker to stop
	quit <- true

	// block until worker confirms it fully stopped.
	// without this: main exits immediately, runtime kills worker mid-cleanup,
	// "Stopping Worker safely..." may never print.
	<-done

	fmt.Println("Main execution completed cleanly.")
}