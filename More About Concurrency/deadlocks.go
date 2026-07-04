package main

import (
	"fmt"
	"sync"
	"time"
)

/*
====================================================
DEADLOCK
====================================================

A deadlock occurs when two or more goroutines
wait forever for each other to release resources.

As a result:

• No goroutine can make progress.
• Program becomes permanently blocked.
• Go runtime eventually panics with:

    fatal error: all goroutines are asleep - deadlock!

----------------------------------------------------

Most Common Causes

1. Channel Deadlocks
2. Mutex Deadlocks
3. WaitGroup Misuse

This file demonstrates Mutex Deadlocks.

====================================================
*/

// ====================================================
// DEADLOCK EXAMPLE
// ====================================================
//
// Lock Order:
//
// Goroutine 1
//
//	mu1 -> mu2
//
// Goroutine 2
//
//	mu2 -> mu1
//
// This inconsistent lock order creates a
// circular wait, leading to a deadlock.
func deadlock() {

	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		// Lock first mutex.
		mu1.Lock()
		fmt.Println("Goroutine 1 locked mu1")

		time.Sleep(1 * time.Second)

		// Waiting for mu2...
		// But Goroutine 2 already owns it.
		mu2.Lock()

		fmt.Println("Goroutine 1 locked mu2")

		mu1.Unlock()
		mu2.Unlock()
	}()

	go func() {
		defer wg.Done()

		// Lock second mutex.
		mu2.Lock()
		fmt.Println("Goroutine 2 locked mu2")

		time.Sleep(1 * time.Second)

		// Waiting for mu1...
		// But Goroutine 1 already owns it.
		mu1.Lock()

		fmt.Println("Goroutine 2 locked mu1")

		mu2.Unlock()
		mu1.Unlock()
	}()

	/*
		Execution Timeline

		Goroutine 1
		    Lock(mu1)
		        ↓
		    Wait(mu2)

		Goroutine 2
		    Lock(mu2)
		        ↓
		    Wait(mu1)

		G1 waits for G2
		G2 waits for G1

		↓

		Circular Wait

		↓

		Deadlock
	*/

	wg.Wait()

	// Never reached.
	fmt.Println("Deadlock Function Completed")
}

// ====================================================
// DEADLOCK PREVENTION
// ====================================================
//
// Always acquire multiple mutexes
// in the SAME order.
//
// Goroutine 1:
//
//	mu1 -> mu2
//
// Goroutine 2:
//
//	mu1 -> mu2
//
// Since everyone follows the same order,
// circular waiting cannot occur.
func noDeadlock() {

	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		mu1.Lock()
		fmt.Println("Goroutine 1 locked mu1")

		time.Sleep(1 * time.Second)

		mu2.Lock()
		fmt.Println("Goroutine 1 locked mu2")

		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() {
		defer wg.Done()

		// Same locking order.
		mu1.Lock()
		fmt.Println("Goroutine 2 locked mu1")

		time.Sleep(1 * time.Second)

		mu2.Lock()
		fmt.Println("Goroutine 2 locked mu2")

		mu2.Unlock()
		mu1.Unlock()
	}()

	wg.Wait()

	fmt.Println("No Deadlock Function Completed")
}

func deadlocks() {

	// Safe Example
	noDeadlock()

	// Deadlock Example
	//
	// Uncomment to observe deadlock.
	//
	// deadlock()

	time.Sleep(4 * time.Second)
}

/*
====================================================
1-MINUTE REVISION
====================================================

Deadlock

Two or more goroutines wait forever
for each other.

----------------------------------------------------

Deadlock Conditions

✓ Mutual Exclusion

✓ Hold and Wait

✓ No Preemption

✓ Circular Wait

----------------------------------------------------

Mutex Deadlock

Goroutine 1

Lock(mu1)
    ↓
Wait(mu2)

Goroutine 2

Lock(mu2)
    ↓
Wait(mu1)

↓

Circular Wait

↓

Deadlock

----------------------------------------------------

Avoiding Deadlocks

✓ Always lock mutexes
  in the same order.

✓ Keep critical sections short.

✓ Unlock immediately after use.

✓ Use defer Unlock() when appropriate.

----------------------------------------------------

Interview Questions

Q. What is a deadlock?

Q. Why does inconsistent lock order
cause deadlocks?

Q. How do you prevent mutex deadlocks?

Q. What are the four Coffman conditions?

====================================================
*/
