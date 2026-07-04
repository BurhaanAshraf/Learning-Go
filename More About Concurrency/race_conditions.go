package main

import (
	"fmt"
	"sync"
)

/*
====================================================
RACE CONDITIONS & MUTEX
====================================================

Race Condition

Occurs when multiple goroutines access the same
variable concurrently and at least one of them
modifies it.

The final result depends on execution timing,
making the program unpredictable.

Example:

	counter++

Although it looks like one operation, internally
it performs three steps:

	1. Read counter
	2. Increment value
	3. Write updated value

If two goroutines perform these steps at the
same time, one update can overwrite the other,
causing a LOST UPDATE.

----------------------------------------------------

Mutex (Mutual Exclusion)

A mutex protects shared resources by allowing
ONLY ONE goroutine into the critical section
at a time.

Workflow:

	Lock()
	    ↓
	Critical Section
	    ↓
	Unlock()

If another goroutine calls Lock() while the
mutex is already locked, it blocks until the
mutex becomes available.

====================================================
*/

type counter struct {

	// Shared resource.
	counter int

	// Protects access to counter.
	mu sync.Mutex
}

// ====================================================
// increment()
// ====================================================
//
// Safely increments the shared counter.
//
// Lock() ensures only one goroutine can modify
// the counter at any given time.
func (c *counter) increment(wg *sync.WaitGroup) {
	defer wg.Done()

	for range 10_000 {

		// Begin Critical Section.
		c.mu.Lock()

		c.counter++

		// End Critical Section.
		c.mu.Unlock()
	}
}

// ====================================================
// incrementBy2()
// ====================================================
//
// Every goroutine accessing the same shared
// resource MUST use the SAME mutex.
//
// Otherwise the mutex cannot prevent races.
func (c *counter) incrementBy2(wg *sync.WaitGroup) {
	defer wg.Done()

	for range 10_000 {

		c.mu.Lock()

		c.counter += 2

		c.mu.Unlock()
	}
}

func race_conditions() {

	counter := counter{counter: 0}

	var wg sync.WaitGroup

	fmt.Println("Starting Main...")

	// ====================================================
	// WAITGROUP
	// ====================================================
	//
	// WaitGroup synchronizes goroutines.
	//
	// Add(n)
	//     Register n goroutines.
	//
	// Done()
	//     Signals one goroutine finished.
	//
	// Wait()
	//     Blocks until counter becomes zero.
	//
	wg.Add(2)

	// ====================================================
	// GOROUTINES
	// ====================================================
	//
	// Methods can be launched directly.
	//
	// No need for:
	//
	// go func() {
	//     counter.increment(&wg)
	// }()
	//
	// Simply:
	//
	go counter.increment(&wg)
	go counter.incrementBy2(&wg)

	wg.Wait()

	fmt.Println("Final Counter Value:", counter.counter)
}

/*
====================================================
TRY THIS YOURSELF
====================================================

1. Remove:

	c.mu.Lock()
	c.mu.Unlock()

from both methods.

2. Run:

	go run main.go

The output may still appear correct.

This DOES NOT mean the code is safe.

----------------------------------------------------

3. Run:

	go run -race main.go

Go's race detector will detect the race
even if the printed output looks correct.

----------------------------------------------------

4. Increase loop count to:

	10_000_000

Run multiple times.

You'll eventually notice the final value
becomes inconsistent because updates are lost.

====================================================
QUICK REVISION
====================================================

Race Condition

Multiple Goroutines
        +
Shared Variable
        +
At Least One Write
        =
Race Condition

----------------------------------------------------

Mutex

Lock()
    ↓
Critical Section
    ↓
Unlock()

----------------------------------------------------

Critical Section

Code that accesses shared data.

Keep it as small as possible.

----------------------------------------------------

Mutex Rules

✓ Lock before accessing shared data.

✓ Unlock immediately after.

✓ Never forget Unlock().

✓ Never copy a Mutex.

✓ Use ONE mutex for ONE shared resource.

----------------------------------------------------

WaitGroup

Add()
	Register goroutine

Done()
	Goroutine finished

Wait()
	Block until all finish

----------------------------------------------------

Race Detector

go run -race main.go

Automatically detects race conditions.

====================================================
INTERVIEW QUESTIONS
====================================================

Q. What is a race condition?

A. Multiple goroutines access shared data
concurrently and at least one writes to it,
making the result depend on timing.

----------------------------------------------------

Q. Why is counter++ not atomic?

A. It performs:

	Read
	Increment
	Write

These steps can interleave.

----------------------------------------------------

Q. What does a Mutex do?

A. It protects shared resources by allowing
only one goroutine into the critical section.

----------------------------------------------------

Q. What is a critical section?

A. Code that accesses shared data.

----------------------------------------------------

Q. Why should a Mutex never be copied?

A. Copying creates independent locks, breaking
synchronization and leading to undefined behavior.

----------------------------------------------------

Q. Why use WaitGroup?

A. To wait until all goroutines finish.

====================================================
*/
