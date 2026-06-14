package main

import (
	"fmt"
	"sync"
)

// ============================================================================
// sync.Mutex — REVISION NOTES
// ============================================================================
//
// WHAT
//   A Mutex (mutual exclusion lock) lets only ONE goroutine access a shared
//   variable at a time. Other goroutines trying to Lock() will wait until
//   the current holder calls Unlock().
//
// WHY
//   Without a mutex, multiple goroutines reading and writing the same
//   variable at the same time causes a DATA RACE — the final value is
//   unpredictable. The Go race detector (`go run -race`) catches this.
//
//   Example of the problem (no mutex):
//     counter++
//   This looks like one operation but is actually three:
//     1. read current value
//     2. add 1
//     3. write new value back
//   Two goroutines can interleave these steps and overwrite each other.
//
// THE 3 METHODS
//   mu.Lock()    — acquires the lock. Blocks if another goroutine holds it.
//   mu.Unlock()  — releases the lock so another goroutine can acquire it.
//   mu.TryLock() — tries to acquire without blocking. Returns true if
//                  it succeeded, false if another goroutine already holds it.
//                  Use when you want to do something else instead of waiting.
//
// ALWAYS USE defer mu.Unlock() RIGHT AFTER Lock()
//   c.mu.Lock()
//   defer c.mu.Unlock()
//   This guarantees Unlock() is called even if the function panics or
//   returns early. Forgetting Unlock() causes every other goroutine to
//   wait forever (deadlock).
//
// ALWAYS PASS MUTEX BY POINTER
//   func increment(c *counter)  ← correct
//   func increment(c counter)   ← wrong: each call gets its own copy of
//                                  the mutex, so locking does nothing.
//
// ============================================================================
// COMMON MISTAKES
// ============================================================================
//
// 1. defer mu.Unlock() inside a loop
//
//    for range 1000 {
//        mu.Lock()
//        defer mu.Unlock()  // BAD: defer runs when the FUNCTION returns,
//        counter++           //      not at the end of each loop iteration.
//    }                       //      Lock is called 1000 times, Unlock runs
//                            //      once at the end — deadlock on iteration 2.
//
//    Fix: call Unlock() directly, not with defer, inside a loop.
//    for range 1000 {
//        mu.Lock()
//        counter++
//        mu.Unlock()         // runs at end of this iteration, not the function
//    }
//
// 2. Forgetting to Unlock
//    Any goroutine that calls Lock() after this will block forever.
//
// 3. Locking twice from the same goroutine
//    sync.Mutex is NOT re-entrant. If a goroutine calls Lock() while it
//    already holds the lock, it deadlocks with itself.
//
// 4. Protecting only the write, not the read
//    Reading a shared variable while another goroutine writes it is also
//    a data race. Reads need to be protected too.
//    ✗  return c.val          // unprotected read
//    ✓  c.mu.Lock(); defer c.mu.Unlock(); return c.val
//
// ============================================================================
// MUTEX AND PERFORMANCE
// ============================================================================
//
// CONTENTION
//   When multiple goroutines try to Lock() at the same time, only one wins.
//   The rest wait. This waiting is called contention. High contention
//   (many goroutines fighting over the same mutex) slows your program down.
//
// GRANULARITY
//   Lock only the smallest piece of code that needs protection.
//   Holding a lock longer than needed increases contention.
//   ✗  mu.Lock(); doSlowNetworkCall(); counter++; mu.Unlock()
//   ✓  result := doSlowNetworkCall(); mu.Lock(); counter++; mu.Unlock()
//
// ============================================================================
// WHEN TO USE MUTEX vs CHANNELS
// ============================================================================
//
//   Use a Mutex when:
//   - You have shared state (a counter, a map, a struct) that multiple
//     goroutines need to read and write.
//   - The access pattern is simple: lock, modify, unlock.
//
//   Use a Channel when:
//   - You want to pass data between goroutines (producer/consumer).
//   - You want to signal events (done, cancel, trigger).
//
// ============================================================================
// INTERVIEW POINTS
// ============================================================================
//
//   Q: What is a data race?
//   A: Two goroutines accessing the same variable at the same time, where
//      at least one is a write. The result is unpredictable.
//
//   Q: How do you detect a race condition?
//   A: go run -race main.go
//
//   Q: Why use defer mu.Unlock()?
//   A: Ensures Unlock() always runs, even on early return or panic.
//      But don't use defer inside a loop — it won't unlock per iteration.
//
//   Q: Is sync.Mutex re-entrant?
//   A: No. If the same goroutine calls Lock() twice without unlocking,
//      it deadlocks with itself.
//
//   Q: Mutex vs Channel — when to use which?
//   A: Mutex for protecting shared state. Channel for communicating
//      between goroutines.
//
// ============================================================================

// ============================================================================
// EXAMPLE 1: Mutex inside a struct (most common real-world pattern)
// ============================================================================
//
// Embedding the mutex in the same struct as the data it protects is the
// standard Go pattern. It keeps the lock and the data it guards together.
// Always group them: mutex first, then the fields it protects.

type counter struct {
	mu  sync.Mutex
	val int // protected by mu — only access inside Lock/Unlock
}

// increment is a standalone function that takes a pointer to counter.
// Passing by pointer is required — if passed by value, the mutex is
// copied and the lock protects nothing.
func increment(c *counter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.val++
}

// getVal shows that READS also need protection.
// Without the lock here, another goroutine could be writing c.val at
// the same time we're reading it — that's a data race.
func (c *counter) getVal() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.val
}

// doIncrement: 10 goroutines each call increment() 1000 times.
// Expected final value: 10,000. Without the mutex, you'd get less
// because goroutines would overwrite each other's writes.
func doIncrement() {
	var wg sync.WaitGroup
	c := &counter{}
	numGoRoutines := 10

	for range numGoRoutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				increment(c)
				// c.val++  ← without the mutex this is a data race.
				//             Run `go run -race` to see it caught.
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Final Val: %d\n", c.getVal()) // should always be 10000
}

// ============================================================================
// EXAMPLE 2: Inline mutex (no struct wrapper)
// ============================================================================
//
// When the shared variable is local to one function, you don't need a struct.
// Declare the mutex alongside the variable and pass both to the goroutine
// via closure.
//
// KEY FIX FROM ORIGINAL CODE:
// The original used `defer mu.Unlock()` inside the for loop — which only
// runs when the whole function returns, not per iteration. That means:
//   iteration 1: Lock() ✓
//   iteration 2: Lock() → DEADLOCKS (still holding lock from iteration 1)
// Fix: call mu.Unlock() directly (no defer) inside the loop.

func doIncrementInline() {
	var (
		count int
		mu    sync.Mutex
		wg    sync.WaitGroup
	)

	numGoRoutines := 5

	increment := func() {
		defer wg.Done()
		for range 1000 {
			mu.Lock()
			count++     // protected section — only one goroutine here at a time
			mu.Unlock() // no defer here — must unlock per iteration, not per function
		}
	}

	wg.Add(numGoRoutines)
	for range numGoRoutines {
		go increment()
	}

	wg.Wait()
	fmt.Printf("Final counter value: %d\n", count) // should always be 5000
}

func Mutex() {
	fmt.Println("=== EXAMPLE 1: Mutex in struct ===")
	doIncrement()

	fmt.Println("\n=== EXAMPLE 2: Inline mutex ===")
	doIncrementInline()
}
