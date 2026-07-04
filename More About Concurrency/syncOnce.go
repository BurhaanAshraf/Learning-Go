package main

import (
	"fmt"
	"sync"
)

/*
================================================================
CORE CONCEPT — sync.Once
================================================================

WHY IT EXISTS:
  Some setup (DB connections, config loading, singleton creation)
  must run EXACTLY ONCE, even when many goroutines might trigger
  it concurrently. sync.Once guarantees this without you having
  to hand-roll the synchronization yourself.

WHY NOT JUST A bool FLAG + MUTEX?
  You could write:
    mu.Lock()
    if !initialized { initialize(); initialized = true }
    mu.Unlock()
  This works, but every single call pays the Lock/Unlock cost
  forever, even long after initialization is done. sync.Once
  internally uses an atomic check first, so after the first call
  it's essentially just an atomic read — much cheaper for code
  that gets called often (e.g. a lazy-init getter called on every
  request).

HOW IT BEHAVES:
  - First goroutine to call Do() runs the function.
  - Every other goroutine calling Do() — whether concurrently or
    much later — blocks until the first call finishes, then
    returns without running anything.
  - WHICH goroutine "wins" and runs first is not deterministic —
    don't rely on goroutine order/id for that.

IMPORTANT GOTCHA — panics:
  If the function passed to Do() panics, sync.Once still marks
  itself as "done". Future Do() calls will NOT retry the function
  — they'll just return silently. So a failed initialization can
  look, to callers, exactly like a successful one. If your setup
  can fail, check for that failure explicitly elsewhere (e.g. a
  separate error variable), don't assume Do() succeeding means
  init succeeded.

COMMON USES: singleton init, config loading, DB connection setup,
  logger init, any global resource that must init exactly once.
================================================================
*/

var once sync.Once

func initialize() {
	fmt.Println("Initialization executed.")
}

func syncOnce() {
	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Println("Goroutine:", id)

			// Only the FIRST goroutine to reach here runs
			// initialize(). Every other goroutine (concurrent
			// or later) just returns immediately — Do() itself
			// blocks them until the first call finishes, so no
			// goroutine can see a "half-initialized" state.
			once.Do(initialize)
		}(i)
	}

	wg.Wait()
	fmt.Println("Main Gracefully Executed...")
}

/*
================================================================
HOW sync.Once WORKS
================================================================
Any goroutine -> once.Do(fn)
                     │
          Has fn already run?
             │             │
            No            Yes
             │             │
        Run fn         Return immediately
             │
      Mark as done
             │
   (all future calls take the "Yes" path)

Note: if two goroutines call Do() at the same instant, one runs
fn while the other BLOCKS (doesn't just skip) until fn finishes —
this guarantees no caller ever proceeds before init is complete.
================================================================
QUICK REVISION
================================================================
sync.Once   -> guarantees one-time execution, thread-safe.
Do(fn)      -> first call runs fn; all later calls are no-ops.
Concurrent calls -> only one runs fn; others block until it's done.
Panics      -> Once still marks itself "done" — no automatic retry.
Cost        -> cheap after first call (atomic check, not a full lock).

Common uses: singleton init, DB connections, config loading,
             logger setup, global resource setup.
================================================================
INTERVIEW Q&A
================================================================
Q: What problem does sync.Once solve?
A: Guarantees one-time execution of initialization code safely
   across concurrent goroutines.

Q: Is sync.Once thread-safe?
A: Yes — multiple goroutines can call Do() simultaneously; only
   one executes the function, others block until it's done.

Q: Does sync.Once execute once per goroutine, or once total?
A: Once total, for the entire Once instance — not per goroutine.

Q: Can Do() run a different function on a later call?
   e.g. once.Do(funcA); once.Do(funcB)
A: No — only the first Do() call's function ever runs. Later
   calls are ignored regardless of what function they pass.

Q: What happens if the function passed to Do() panics?
A: Once is still marked "done" — it will not retry on subsequent
   Do() calls, even though initialization effectively failed.
   This is a common source of subtle bugs.

Q: Why prefer sync.Once over a bool flag guarded by a Mutex?
A: Same correctness, but sync.Once is cheaper on the common path
   after the first call — no full lock/unlock needed once done,
   just an atomic check.
================================================================
*/
