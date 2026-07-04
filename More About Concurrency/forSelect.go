package main

import (
	"fmt"
	"time"
)

/*
================================================================
CORE CONCEPT — select
================================================================

WHAT select DOES:
  select waits on MULTIPLE channels at once and runs whichever
  case becomes ready FIRST. If more than one is ready simultaneously,
  Go picks one at random (not top-to-bottom) — never assume order.

WHY IT'S NEEDED HERE:
  A goroutine can only block on ONE channel receive at a time
  with plain `<-ch`. But this program needs to react to TWO
  independent events — "tick happened" OR "shutdown requested" —
  without knowing in advance which comes first. select is Go's
  only way to listen to several channels concurrently.

  for { select { ... } } is the standard EVENT LOOP pattern:
  loop forever, block until any one channel fires, handle it,
  repeat. This is how most long-running Go services structure
  their main work loop.

================================================================
time.Ticker
================================================================
  ticker := time.NewTicker(interval)
  Fires by sending the current time into ticker.C every interval.
  ticker.C is receive-only — you only ever read from it.

  Always defer ticker.Stop() — otherwise the internal timer keeps
  running and leaking resources even after you stop using it.
  NOTE: Stop() does NOT close ticker.C — so `case <-ticker.C`
  after Stop() would just block forever, not panic or fire zero
  values. That's fine here since we return right after, but worth
  knowing if you reuse a ticker across a longer-lived program.

================================================================
close(quit) as a shutdown signal
================================================================
  Closing a channel BROADCASTS to every goroutine waiting on it —
  unlike sending a value (which only one receiver gets), a closed
  channel unblocks ALL current and future receives immediately.
  This is why `close()` (not a sent value) is the idiomatic Go
  shutdown signal when multiple goroutines need to hear about it.
================================================================
*/

func Select() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	quit := make(chan struct{})

	// Simulates another goroutine deciding to stop the app
	// after 5 seconds — close() broadcasts the shutdown signal.
	go func() {
		time.Sleep(5 * time.Second)
		close(quit)
	}()

	// EVENT LOOP: block until either channel is ready.
	for {
		select {
		case <-ticker.C: // fires every second
			fmt.Println("Received Tick...")
		case <-quit: // fires once, when quit is closed
			fmt.Println("Quitting...")
			return
		}
	}
}

/*
================================================================
PROGRAM FLOW
================================================================
             ┌── ticker.C ready → print tick, loop again
select ──────┤
             └── quit closed  → print "Quitting...", return
================================================================
QUICK REVISION
================================================================
select        -> waits on multiple channels, runs whichever
                 fires first (random pick if several are ready
                 at once — never assume ordering).
ticker.C      -> receive-only channel, fires every interval.
ticker.Stop() -> stops the timer; does NOT close ticker.C.
close(quit)   -> broadcasts to ALL receivers at once — the
                 idiomatic multi-goroutine shutdown signal.

Modern alternative worth knowing: many newer codebases use
context.Context + ctx.Done() instead of a hand-rolled quit
channel — same broadcast idea, but standardized and composable
with timeouts/cancellation up and down a call chain.
================================================================
INTERVIEW Q&A
================================================================
Q: Timer vs Ticker?
A: Timer fires once. Ticker fires repeatedly at fixed intervals.

Q: Why is select needed instead of just `<-ticker.C`?
A: A plain receive blocks on ONE channel only. select lets a
   goroutine react to whichever of several channels is ready
   first — required here since we listen for both ticks and quit.

Q: What happens if two select cases are ready at the same time?
A: Go picks one at random — don't rely on any particular order.

Q: Why call ticker.Stop()?
A: Releases the ticker's internal timer resources; without it,
   the timer keeps running even if nothing reads ticker.C anymore.

Q: Why close(quit) instead of sending a value?
A: close() broadcasts to every current and future receiver at
   once; a sent value is only ever received by one goroutine.
================================================================
*/
