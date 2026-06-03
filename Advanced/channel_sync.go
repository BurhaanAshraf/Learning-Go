package main

import (
	"fmt"
	"time"
)

/*
=============================================================================
CHANNEL SYNCHRONIZATION — COMPLETE REVISION NOTES
=============================================================================

─────────────────────────────────────────────────────────────────────────────
WHAT IS CHANNEL SYNCHRONIZATION?
─────────────────────────────────────────────────────────────────────────────

Channel synchronization means using a channel as a meeting point between
the main goroutine and worker goroutines so that main waits for a goroutine
to finish before moving on.

The core problem it solves:
    Main does NOT wait for goroutines by default.
    If main exits, ALL goroutines are killed instantly — finished or not.

Without sync:
    main    ----exits----
    worker      ----still running---- (killed, work lost)

With channel sync:
    main    ----blocks here----agrees----continues----
    worker  ----working-----------done----signals----

They meet at the channel. Main cannot pass that point until the worker
signals. That meeting point IS the synchronization.

So channel sync = using a channel as the meeting point where two goroutines
agree before either moves on.

─────────────────────────────────────────────────────────────────────────────
WHY UNBUFFERED CHANNELS FOR SYNC?
─────────────────────────────────────────────────────────────────────────────

Unbuffered channels are the natural fit for synchronization because their
blocking behavior IS the sync mechanism.

With unbuffered channel:
    → Send BLOCKS until a receiver is ready
    → Receive BLOCKS until a sender is ready
    → Both sides MUST meet at the same time (rendezvous)
    → Neither side can skip past it — this is what forces the sync

With buffered channel:
    → Sender just drops value in buffer and moves on immediately
    → Does not need a receiver to be ready
    → You lose the guaranteed rendezvous that makes sync reliable

RULE:
    Unbuffered channels for sync.
    Buffered channels for data flow and decoupling.

─────────────────────────────────────────────────────────────────────────────
THE DONE CHANNEL PATTERN — PUREST FORM OF CHANNEL SYNC
─────────────────────────────────────────────────────────────────────────────

The done channel pattern is the standard way to sync a single goroutine
with main. The idea is simple:

    1. Create an unbuffered channel called done
    2. Launch goroutine — it does its work, sends a signal at the very end
    3. Main blocks on receive — cannot move on until signal arrives
    4. Once signal arrives, main unblocks and continues

The VALUE sent does not matter at all. Nobody reads it.
The ARRIVAL of the value is the signal — it means "I am done."
That is why we send true, 0, 1 — anything, it does not matter.

CORRECT DIRECTION:
    Goroutine sends last → main receives last ✅
    Main waits for goroutine to finish all its work before signaling.

WRONG DIRECTION:
    Main sends → goroutine receives last ⚠️
    Main only waits for goroutine to be READY to receive, not to finish.
    Whatever goroutine does after receiving is not waited on by main.

─────────────────────────────────────────────────────────────────────────────
SLEEP IS NOT SYNC
─────────────────────────────────────────────────────────────────────────────

A common mistake is using time.Sleep on main to "wait" for goroutines:

    go func() { ... }()
    time.Sleep(2 * time.Second)  // hoping goroutine finishes in 2 seconds
    fmt.Println("done")

This is NOT sync. This is guessing. Main does not know or care what
goroutines are doing during the sleep. It just wakes up after 2 seconds
and moves on regardless of whether goroutines finished.

Real sync is when main itself is blocked on a channel receive — because
then main CANNOT move on until the goroutine actually sends the signal.
No guessing, no timing, just a hard block until signal arrives.

RULE:
    If the receive is not on main, main is not syncing — it is waiting blindly.

─────────────────────────────────────────────────────────────────────────────
SYNCING MULTIPLE GOROUTINES
─────────────────────────────────────────────────────────────────────────────

When you have N goroutines, you need N signals — one from each goroutine.
Main must receive N times, once per goroutine, before moving on.

Use a buffered channel with capacity N so goroutines do not block each
other when sending their signals. Main still blocks on each receive.

RULE:
    Receives must be on MAIN directly — not inside another goroutine.
    If receives are inside a goroutine, main is not actually syncing.
    Main will just continue past that goroutine and exit.

─────────────────────────────────────────────────────────────────────────────
FOR LOOP vs INDIVIDUAL GOROUTINES
─────────────────────────────────────────────────────────────────────────────

These two are identical in behavior:

    with for loop
    for i := range 3 {
        go func(id int) { done <- id }(i)
    }
    for range 3 { <-done }

    individually
    go func() { done <- 0 }()
    go func() { done <- 1 }()
    go func() { done <- 2 }()
    <-done
    <-done
    <-done

The for loop is just cleaner and scalable. Imagine doing it individually
for 100 goroutines. The underlying mechanism is exactly the same.

─────────────────────────────────────────────────────────────────────────────
LOOP VARIABLE CAPTURE BUG
─────────────────────────────────────────────────────────────────────────────

When launching goroutines inside a loop, always pass the loop variable
as a parameter to the goroutine function. Never use it directly inside.

    WRONG:
        for i := range 3 {
            go func() {
                fmt.Println(i)  // i may have changed by the time this runs
            }()
        }

    CORRECT:
        for i := range 3 {
            go func(id int) {
                fmt.Println(id)  // id is a copy, safe to use
            }(i)
        }

Since goroutines run concurrently, by the time a goroutine reads i,
the loop may have already incremented it. Passing it as a parameter
creates a copy that belongs to that goroutine alone.

─────────────────────────────────────────────────────────────────────────────
EQUAL SENDERS AND RECEIVERS
─────────────────────────────────────────────────────────────────────────────

For every send there must be a matching receive, and for every receive
there must be a matching send. If they are unequal, the unmatched side
blocks forever with nobody to unblock it → deadlock.

    3 goroutines sending → main must receive 3 times
    1 goroutine sending  → main must receive 1 time

─────────────────────────────────────────────────────────────────────────────
GOROUTINE BLOCKED ON RECEIVE WITH NO VALUE — IS IT A PROBLEM?
─────────────────────────────────────────────────────────────────────────────

A goroutine sitting blocked on a receive is NOT a problem by itself.
Go does not care about blocked goroutines as long as main can still
make progress.

It only becomes a deadlock when EVERY goroutine including main is
blocked with no way out — nobody left to send or receive.

    go func() {
        <-done  // blocked, waiting — this alone is fine
    }()
    fmt.Println("main continues normally")  // main is not blocked, no problem

    vs

    go func() {
        <-done  // blocked
    }()
    <-done      // main also blocked, nobody sends → DEADLOCK

─────────────────────────────────────────────────────────────────────────────
SUMMARY — KEY RULES TO REMEMBER
─────────────────────────────────────────────────────────────────────────────

1. Main does not wait for goroutines by default — channel sync fixes this.

2. Channel sync = using a channel as a meeting point between goroutines.

3. Use UNBUFFERED channels for sync — their blocking behavior is the sync.

4. Correct pattern: goroutine does work → sends signal last → main receives.

5. Sleep is not sync — it is just guessing with time.

6. Receives must be on MAIN directly for real sync — not inside a goroutine.

7. Every send needs a matching receive and vice versa — unmatched = deadlock.

8. A blocked goroutine is fine as long as main can still make progress.

9. Loop variable must be passed as parameter to goroutine — never use directly.

10. For loop goroutines and individual goroutines are identical in behavior —
    for loop is just cleaner and scalable.

11. You do not need to close every channel. Closing is only a signal to
    receivers (like for range) that no more data is coming. If no one is
    ranging over it, leaving a channel unclosed does NOT cause a memory leak;
    the Go garbage collector will clean it up naturally.

12. A buffered channel acts like an unbuffered channel once it is full.
    If you have 3 workers but a channel capacity of 1, the first worker
    drops its signal and exits immediately. Because the buffer is now full,
    the channel behaves like an unbuffered channel for the remaining 2
    workers—they will block on their send (`done <- id`) and can only finish
    one by one as main actively receives values and frees up space.

=============================================================================
*/

// =============================================================================
// EXAMPLE 1 — BASIC CHANNEL SYNC (DONE CHANNEL PATTERN)
// =============================================================================
//
// Demonstrates: the purest form of channel sync.
// Goroutine does its work, sends a signal at the end.
// Main blocks on receive until signal arrives, then continues.
// The value sent does not matter — only the arrival matters.

func basicSync() {

	done := make(chan bool) // unbuffered — forces rendezvous, perfect for sync

	go func() {
		fmt.Println("Working...")
		time.Sleep(2 * time.Second) // simulating work
		done <- true                // work is done, send signal — value does not matter
	}()

	// main blocks here — cannot move past this line until goroutine sends signal.
	// once signal arrives, main unblocks and continues.
	<-done

	fmt.Println("Finished...")
}

// =============================================================================
// EXAMPLE 2 — BASIC DATA EXCHANGE VIA CHANNEL
// =============================================================================
//
// Demonstrates: sending an actual value through a channel, not just a signal.
// Here the value itself matters, not just the arrival.
// Both sides block until the other is ready — unbuffered rendezvous.

func dataExchange() {

	channel := make(chan int) // unbuffered — sender and receiver must meet

	go func() {
		channel <- 9            // blocks here until main is ready to receive
		fmt.Println("Value Sent")
	}()

	// blocks here until goroutine sends the value.
	// once value arrives, main unblocks, receives 9, and continues.
	value := <-channel
	fmt.Println("Value Received", value)
}

// =============================================================================
// EXAMPLE 3 — SYNCING MULTIPLE GOROUTINES (CORRECTED VERSION)
// =============================================================================
//
// Demonstrates: waiting for N goroutines to finish before main continues.
// Uses a buffered channel with capacity N as a performance optimization.
// This allows goroutines to drop their signal and exit immediately, rather 
// than blocking and waiting for main to iterate through the loop to receive it.
// Receives are on MAIN directly — this is what makes it real sync.
//
// Common mistakes fixed from original code:
//   1. Receives moved from goroutine to main directly
//   2. Loop variable passed as parameter (id) to avoid capture bug
//   3. No sleep guessing — main blocks on actual channel receives

func syncMultipleGoroutines() {

	numGoroutines := 3
	done := make(chan int, numGoroutines) // buffered with capacity N so goroutines
	                                      // do not block each other on send

	for i := range numGoroutines {
		go func(id int) { // pass i as id — safe copy for this goroutine
			fmt.Printf("Goroutine %d is working...\n", id)
			time.Sleep(time.Second)
			done <- id // signal completion — value is id but could be anything
		}(i)
	}

	// receive directly on main — NOT inside a goroutine.
	// main blocks here for each receive, cannot move on until all 3 goroutines send.
	// this is real sync — main is blocked, not sleeping.
	for range numGoroutines {
		<-done
	}

	fmt.Println("All Goroutines Are Finished...")
}

// =============================================================================
// EXAMPLE 4 — SYNCING MULTIPLE GOROUTINES (INDIVIDUALLY, NO LOOP)
// =============================================================================
//
// Demonstrates: exact same behavior as example 3 but without a for loop.
// Proves that for loop and individual goroutines are identical in behavior.
// For loop is just cleaner and scalable for larger numbers.

func syncMultipleGoroutinesIndividual() {

	done := make(chan int, 3) // buffered with capacity 3

	// three individual goroutines — same as launching with a for loop
	go func() {
		fmt.Println("Goroutine 0 is working...")
		time.Sleep(time.Second)
		done <- 0
	}()

	go func() {
		fmt.Println("Goroutine 1 is working...")
		time.Sleep(time.Second)
		done <- 1
	}()

	go func() {
		fmt.Println("Goroutine 2 is working...")
		time.Sleep(time.Second)
		done <- 2
	}()

	// three receives on main — one per goroutine.
	// main blocks until all three have sent their signal.
	<-done
	<-done
	<-done

	fmt.Println("All Goroutines Are Finished...")
}

// =============================================================================
// EXAMPLE 5 — SYNCHRONISING DATA EXCHANGE WITH FOR RANGE
// =============================================================================
//
// Demonstrates: sending a stream of values through a channel and receiving
// them using for range. close(ch) is essential here — it signals for range
// to stop looping once all values are drained.
//
// for range on a channel:
//   → actively pulls values out of the channel one by one on each iteration
//   → it IS a receive operation — same as value := <-ch
//   → blocks until next value arrives
//   → stops automatically when channel is closed and empty
//   → if channel is not closed, for range blocks forever → deadlock
//
// Three ways to receive from a channel:
//   value := <-ch              → receive once
//   for value := range ch {}   → receive until channel is closed
//   value, ok := <-ch          → receive and check if channel is still open

func syncDataExchange() {

	data := make(chan string)

	go func() {
		for i := range 5 {
			// fmt.Sprintf is the idiomatic way to build strings with numbers
			data <- fmt.Sprintf("Hello %d", i)
			time.Sleep(100 * time.Millisecond)
		}
		// close signals for range to stop once all values are drained.
		// without close, for range would block forever after last value → deadlock.
		// channels follow FIFO order so values arrive in the order they were sent.
		close(data)
	}()

	// for range receives values one by one, blocks between each until next arrives.
	// stops automatically once channel is closed and all values are drained.
	for value := range data {
		fmt.Println(value)
	}
}

// =============================================================================
// EXAMPLE 6 — WHAT HAPPENS WITHOUT close() — INTENTIONAL DEADLOCK DEMO
// =============================================================================
//
// Demonstrates: why close(ch) is required when using for range.
// Without close, for range has no way to know if more values are coming.
// It blocks forever waiting for the next value that never arrives → deadlock.
//
// Uncomment to see: fatal error: all goroutines are asleep — deadlock!

func deadlockWithoutClose() {

	data := make(chan string)

	go func() {
		for i := range 5 {
			data <- fmt.Sprintf("Hello %d", i)
		}
		 close(data) 
		//  no close(data) here — for range on main will block forever
		//  after receiving all 5 values, waiting for a 6th that never comes
	}()

	for value := range data { // blocks forever after 5th value → deadlock
		fmt.Println(value)
	}
}

func channelSync() {
	fmt.Println("=== Example 1: Basic Sync ===")
	basicSync()

	fmt.Println("\n=== Example 2: Data Exchange ===")
	dataExchange()

	fmt.Println("\n=== Example 3: Sync Multiple Goroutines (Loop) ===")
	syncMultipleGoroutines()

	fmt.Println("\n=== Example 4: Sync Multiple Goroutines (Individual) ===")
	syncMultipleGoroutinesIndividual()

	fmt.Println("\n=== Example 5: Sync Data Exchange ===")
	syncDataExchange()

	fmt.Println("\n=== Example 6: Deadlock Wihtout Close ===")
	deadlockWithoutClose()
}