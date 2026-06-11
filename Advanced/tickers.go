package main

import (
	"context"
	"fmt"
	"time"
)

/*
===============================================================================
STUDY NOTE 1: THE FOUNDATIONAL MECHANICS OF TICKERS VS TIMERS
===============================================================================
- time.Timer: Fires EXACTLY ONCE after a duration. Used for timeouts/delays.
- time.NewTicker: Fires REPEATEDLY at a set interval. Used for periodic loops.
- CRITICAL MECH: Both types deliver data via an underlying channel (.C).
- LEAK WARNING: Tickers interact directly with the OS runtime scheduler. If you
  do not explicitly call .Stop(), the ticker remains alive in memory, leaking
  CPU cycles indefinitely, even if the function containing it returns!
===============================================================================
*/

// Tickersexe demonstrates basic ticker consumption via range and scope isolation.
func Tickersexe() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // REVISION RULE: Always defer Stop() immediately after creation!

	// -------------------------------------------------------------------------
	// PITFALL WRONG: Declaring state variables in the outer routine scope when
	// a background goroutine modifies them can lead to Silent Data Races.
	// REFIX: Keep the state variable ('i') inside the goroutine's scope
	// if the main routine doesn't need to read it concurrently.
	// -------------------------------------------------------------------------
	done := make(chan struct{})

	go func() {
		i := 1 // SAFE: Isolated completely inside this goroutine's stack memory frame

		// MECHANIC: 'for range ticker.C' reads from the channel sequentially.
		// It blocks automatically until the next tick arrives.
		for range ticker.C {
			i *= 2
			if i > 20 {
				break // Exits the loop, allowing execution to reach close(done)
			}
			fmt.Println("Tickersexe item:", i)
		}

		// ---------------------------------------------------------------------
		// PITFALL WRONG: Forgetting to close the signaling channel on exit will
		// cause a "fatal error: all goroutines are asleep - deadlock!" in main.
		// ---------------------------------------------------------------------
		close(done)
	}()

	<-done // BLOCKING POINT: Intentionally holds the main thread until the goroutine finishes.
}

// periodicTask simulates a repetitive background worker function.
func periodicTask() {
	fmt.Println("Performing Periodic Task at:", time.Now().Format("15:04:05"))
}

// runPeriodicTask demonstrates explicit loop execution control using select blocks.
func runPeriodicTask() {
	newticker := time.NewTicker(1 * time.Second)
	defer newticker.Stop()

	done := make(chan struct{})

	go func() {
		i := 0
		for {
			select {
			case <-newticker.C:
				periodicTask()
				i++
				if i == 5 {
					// ---------------------------------------------------------
					// REVISION NOTE: SIGNAL VS EXIT
					// close(done) unblocks the main thread.
					// 'return' terminates this background goroutine.
					// PITFALL: If you omit 'return', the 'for {}' loop spins
					// forever in the background, causing a "Goroutine Leak".
					// ---------------------------------------------------------
					close(done)
					return
				}
			}
		}
	}()

	<-done // Blocks main until the 5 periodic tasks complete.
	fmt.Println("We are done with our periodic checks...")
}

// stoppingGracefully shows basic coordinate timing multiplexing using standard select.
func stoppingGracefully() {
	ticker := time.NewTicker(1 * time.Second)
	stopper := time.After(5 * time.Second) // MECHANIC: time.After returns a <-chan time.Time that fires once.
	defer ticker.Stop()

	for {
		select {
		case <-stopper:
			// REVISION NOTE: time.After does not need manual stopping because it self-expires.
			fmt.Println("stoppingGracefully: Ticker stopped gracefully...")
			return // Drops out of the function completely

		case <-ticker.C:
			fmt.Println("Periodic Routine Checking...", time.Now().Format("15:04:05"))
		}
	}
}

// multipleTickers showcases dynamic ticker manipulation and the Nil-Channel trick.
func multipleTickers() {
	ticker1 := time.NewTicker(1 * time.Second)
	ticker2 := time.NewTicker(5 * time.Second)
	totalTime := time.NewTimer(10 * time.Second)

	// Always clean up resources regardless of exit pathway
	defer totalTime.Stop()
	defer ticker1.Stop()
	defer ticker2.Stop()

	// -------------------------------------------------------------------------
	// CRITICAL PITFALL: THE CHANNELS ARE READ-ONLY STRUCT FIELDS!
	// You cannot execute `ticker1.C = nil` directly (Compilation Error).
	// SOLUTION: Extract the references into local variables which CAN be reassigned.
	// -------------------------------------------------------------------------
	ch1 := ticker1.C
	ch2 := ticker2.C

	for {
		select {
		case <-ch1:
			fmt.Println("Executing task 1...")

		case <-ch2:
			fmt.Println("Executing Task 2 -> Modifying Ticker Frequencies")

			ticker1.Stop() // Stops the internal OS timer subsystem...

			// -----------------------------------------------------------------
			// CRITICAL MECHANIC: THE NIL CHANNEL TRICK
			// Calling ticker1.Stop() DOES NOT close or empty the channel.
			// If a tick was buffered right before stopping, 'case <-ch1:' can still fire!
			// Reassigning ch1 to 'nil' forces the Go runtime select statement
			// to completely ignore this case block on all subsequent iterations.
			// -----------------------------------------------------------------
			ch1 = nil

			// Dynamic manipulation: accelerate ticker2 to sample data faster
			ticker2.Reset(1 * time.Second)

		case <-totalTime.C:
			fmt.Println("Stopping Ticker 2 via Global Timeout Control")
			fmt.Println("Exiting Periodic Check Gracefully...")
			return
		}
	}
}

func ticker() {
	fmt.Println("=== STARTING REVISION PLAYGROUND ===")

	// -------------------------------------------------------------------------
	// QUICK QUICK WORKER NOTE: THE ONE-TICK DISPATCH
	// -------------------------------------------------------------------------
	ticker := time.NewTicker(1 * time.Second)
	for tick := range ticker.C {
		fmt.Println("Initial One-Off Test Tick Time:", tick.Format("15:04:05"))
		break // Exits the range loop immediately after 1 iteration.
	}
	ticker.Stop() // Manual stop needed because we broke out of the loop early!

	// Run Sub-Experiments
	Tickersexe()
	runPeriodicTask()
	stoppingGracefully()

	// -------------------------------------------------------------------------
	// REVISION NOTE: IDIOMATIC CHANNEL SIGNALS
	// PITFALL: Beginners often write `done <- struct{}{}` followed by `close(done)`.
	// This is redundant. Reading from a closed channel instantly returns the
	// zero-value without blocking. Therefore, simply calling `close(done)` is
	// sufficient to unblock any thread listening with `<-done`.
	// -------------------------------------------------------------------------
	done := make(chan struct{})
	go func() {
		multipleTickers()
		close(done) // Instantly unblocks the main thread
	}()

	<-done

	// -------------------------------------------------------------------------
	// BONUS ADVANCED NOTE: UPGRADING FROM RAW TIMERS TO CONTEXT
	// Production-grade Go microservices rarely use hardcoded timeouts inside functions.
	// Instead, they receive a `context.Context` from upstream (HTTP requests, RPCs).
	// Below is the idiomatic pattern showing how to use context cancellation.
	// -------------------------------------------------------------------------
	fmt.Println("\n=== BONUS: CONTEXT DRIVEN REFACTOR ===")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // Avoids context leaks

	contextWorker(ctx)

	fmt.Println("Execution Completed Safely.")
}

// contextWorker models a production loop listening to a context deadline.
func contextWorker(ctx context.Context) {
	ctxTicker := time.NewTicker(500 * time.Millisecond)
	defer ctxTicker.Stop()

	for {
		select {
		case <-ctx.Done(): // Fires automatically when the 3-second context timeout expires
			fmt.Println("Context Terminated Lifecycle Loop:", ctx.Err())
			return
		case <-ctxTicker.C:
			fmt.Println("Context-managed tick running...")
		}
	}
}
