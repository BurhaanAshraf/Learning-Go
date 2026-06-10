package main

import (
	"fmt"
	"time"
)

// ============================================================================
// GO TIME PRIMITIVES — ARCHITECTURAL & INTERVIEW REVISION NOTES
// ============================================================================
//
// 1. THE CORE ENGINE: HOW TIMERS WORK UNDER THE HOOD
//    When you execute time.NewTimer(), Go does NOT spawn a dedicated OS thread
//    or a lightweight goroutine to count down. Doing so would destroy CPU performance.
//    Instead, Go manages time primitives via a highly optimized scheduling layer:
//
//    - THE QUATERNARY MIN-HEAP: Go maintains a balanced min-heap structure internally
//      (attached directly to Go's runtime logical processors, or 'P' structs).
//      Timers are registered here, sorted strictly by absolute expiration timestamp.
//    - THE RUNTIME POLLER: A centralized, highly optimized background runtime thread
//      monitors this heap. When the clock passes the root node's timestamp, it takes
//      the value and pushes it into the `Timer.C` channel.
//    - THE BUFFER GUARANTEE: The internal channel `C` is created with a buffer size
//      of exactly 1 (`make(chan time.Time, 1)`). This ensures that when the runtime
//      writes the timestamp, it does NOT block the entire runtime engine if your
//      worker goroutine isn't actively listening yet.
//
// 2. THE HIGH-TRAFFIC MEMORY LEAK TRAP (time.After vs time.NewTimer)
//    - time.NewTimer(d) allocates a struct containing control handles (.Stop() and .Reset())
//      along with the channel C. It returns a pointer, allowing you to manage its lifecycle.
//    - time.After(d) is a wrapper. It creates a timer internally, throws away the struct
//      pointer completely, and returns ONLY the read-only channel field (`<-chan time.Time`).
//
//    CRITICAL INTERVIEW TRAP: Because time.After throws away the pointer, you CANNOT call
//    .Stop() on it. If you use `time.After(30 * time.Second)` inside a high-frequency loop
//    or an HTTP handler processing 5,000 requests/sec, and those operations finish fast in
//    10ms, the thrown-away timer structs REMAIN ALIVE and pinned in the runtime heap for
//    the full 30 seconds. They cannot be garbage collected early, causing a massive memory spike.
//
//    RULE: Never use time.After inside loops or repetitive handlers. Allocate a single
//          time.NewTimer outside, manage it with Stop(), and modify it via Reset().
// ============================================================================

// Timeout demonstrates the precise lifecycle of a Timer: Stop, Read, and Reset.
func Timeout() {
	fmt.Println("Starting App")
	
	// Non-blocking allocation. Timer starts immediately in the background runtime heap.
	timer := time.NewTimer(2 * time.Second) 
	fmt.Println("Waiting For Timer.C")

	// --- THE STOP() RACE CONDITION ---
	// timer.Stop() prevents a timer from firing if it hasn't expired yet.
	// RETURNS TRUE: The timer was found and removed before it could expire. 
	//               The channel `timer.C` is guaranteed to be completely clean and empty.
	// RETURNS FALSE: The timer has ALREADY expired and pushed its value to `timer.C`, 
	//                OR it was already stopped previously.
	stopped := timer.Stop() 
	
	if stopped {
		fmt.Println("Timer Stopped")
	} else {
		// DANGER ZONE: If stopped is false, the timestamp value has already crossed the 
		// boundary and is sitting in the channel buffer. We MUST drain it here, otherwise 
		// a future read will consume this stale value instantly instead of waiting!
		<-timer.C 
	} // blocking in nature if we hit the else path without a guaranteed value

	fmt.Println("timer expired")

	// --- THE IDIOMATIC RESET PATTERN ---
	// TRAP: Reset() should only be invoked on stopped or expired timers with drained channels.
	// If you call Reset() without ensuring the channel buffer is empty, the very next read 
	// from `<-timer.C` returns the old, stale timestamp instantly without waiting.
	fmt.Println("Timer Reset")
	
	// Defensively ensuring channel is clear before resetting
	if !timer.Stop() {
		select {
		case <-timer.C: // Successfully drained old stale value out of the buffer
		default:        // Channel was already empty
		}
	}
	
	timer.Reset(3 * time.Second)
	<-timer.C // Blocks execution paths cleanly for 3 seconds
	fmt.Println("After 3 seconds it stopped")
}

func longRunningOperation(stopChan chan struct{}) {
	// PRACTICAL NOTE: We passed a stop channel to make this routine abortable.
	// Without it, the goroutine would continue running 20 seconds in the background
	// even after the calling select statement timed out! (Classic Goroutine Leak).
	for i := range 20 {
		select {
		case <-stopChan:
			fmt.Printf("[Worker] Aborting early at iteration %d\n", i)
			return
		default:
			fmt.Println(i)
			time.Sleep(200 * time.Millisecond) // Sped up for practical testing
		}
	}
}

// processLongOperation demonstrates asynchronous boundary tracking using select.
func processLongOperation() {
	// time.After is perfectly acceptable for one-off select timeouts outside of loops.
	timeout := time.After(1200 * time.Millisecond)
	
	done := make(chan struct{})
	stopWorker := make(chan struct{}) // Anti-leak synchronization handshake channel
	
	go func() {
		longRunningOperation(stopWorker)
		done <- struct{}{}
	}()
	
	// --- THE INVISIBLE GOROUTINE LEAK TRAP ---
	// When a `select` case finishes, the other unchosen case branches are NOT cancelled or killed.
	// If `timeout` fires first, this function exits, but the worker goroutine stays alive in memory 
	// unless explicitly signaled to stop via a close or context broadcast channel.
	select {
	case <-timeout:
		stopWorker <- struct{}{}
		close(stopWorker) // Kills the background goroutine explicitly
		fmt.Println("Operation Timeout - Signalling worker to terminate.")
		<-done            // Wait for worker to exit cleanly before leaving function scope
	case <-done:
		fmt.Println("Operation Completed")
	}
}

// delayedOperations demonstrates firing detached asynchronous tasks.
func delayedOperations() {
	done := make(chan struct{})
	timer := time.NewTimer(2 * time.Second) // Non-blocking allocation
	defer timer.Stop()                      // ALWAYS defer Stop() to clean up heap registrations early

	go func() {
		<-timer.C // Blocks locally in the background thread until the runtime poller pushes a tick
		fmt.Println("Delayed Operation Executed")
		done <- struct{}{}
		close(done)
	}()
	
	fmt.Println("Waiting...")
	<- done // Blocking the main thread, until we receive the value from goroutine
	// time.Sleep(3 * time.Second) // Blocking the main execution thread completely
	fmt.Println("End of Program!")
}

// multipleTimers demonstrates multiplexing multiple time streams concurrently.
func multipleTimers() {
	timer1 := time.NewTimer(1 * time.Second)
	timer2 := time.NewTimer(2 * time.Second)
	totalTime := time.After(3 * time.Second)

	defer timer1.Stop()
	defer timer2.Stop()

	// PRACTICAL NOTE: Used a synchronization handshake instead of letting main blindly sleep.
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			select {
			case <-totalTime:
				fmt.Println("Operation Timed Out - Breaking Multiplexer Loop")
				return // Explicit return terminates the infinite loop and kills the goroutine safely
			case <-timer2.C:
				fmt.Println("Timer 2 Expired...")
			case <-timer1.C:
				fmt.Println("Timer 1 Expired...")
			}
		}
	}()

	<-done // Coordinate cleanly with main thread execution lifecycle
}

func Timer() {
	fmt.Println("=== RUNNING SCENARIO 1: TIMEOUT & RESET LIFE ===")
	Timeout()

	fmt.Println("\n=== RUNNING SCENARIO 2: LEAK SENSITIVE TIMEOUT MANAGEMENT ===")
	processLongOperation()

	fmt.Println("\n=== RUNNING SCENARIO 3: DELAYED ASYNC CHECKS ===")
	delayedOperations()

	fmt.Println("\n=== RUNNING SCENARIO 4: MULTIPLEXING MULTIPLE TIMERS ===")
	multipleTimers()

	fmt.Println("\nEnded...")
}

/* ============================================================================
   HIGH-FREQUENCY INTERVIEW Q&A CARD SUMMARY
   ============================================================================
   
   Q: Does calling timer.Stop() close the underlying timer.C channel?
   A: No. It only prevents the background scheduler from pushing a timestamp to it. 
      It never closes the channel. Closing a channel in Go means data can never be 
      written to it again. Since a timer can be revived via .Reset(), closing the 
      channel would make the object permanently broken.

   Q: Why does a second read from a channel returned by time.After block permanently?
   A: Because the channel has a buffer size of 1 and is never closed. The runtime 
      writes a single value to it upon expiration. The first read drains that value. 
      Since no second write will ever occur, a second read locks the current goroutine 
      into a perpetual deadlock state.

   Q: How does the Go runtime scale when managing millions of active timers?
   A: It scales per runtime logical processor (P), not per OS thread. Go uses a specialized 
      balanced 4-ary (quaternary) min-heap structure which reduces heap lock contention and 
      keeps timer search, insertion, and eviction speeds strictly bounded at O(log N).
============================================================================ */