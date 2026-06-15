package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ============================================================================
// sync/atomic — REVISION NOTES
// ============================================================================
//
// WHAT
//   Atomic operations let you read and write shared variables safely across
//   goroutines — without manually locking and unlocking a mutex.
//   The sync/atomic package provides these operations.
//
// WHY
//   counter++ looks like one step but is actually three:
//     1. read current value
//     2. add 1
//     3. write new value back
//   Two goroutines can interleave these steps and overwrite each other.
//   This is a data race. atomic.AddInt64() does all three steps as one
//   uninterruptible operation — no other goroutine can sneak in between.
//
// ATOMIC vs MUTEX
//   Both solve data races. The difference is scope:
//
//   Atomic  → for simple operations on a single value (increment, compare,
//              load, store). Faster, no manual Lock/Unlock.
//   Mutex   → for protecting a block of code or multiple variables together.
//              More flexible but slightly more overhead.
//
//   Rule of thumb:
//     Single variable, simple operation → use atomic.
//     Multiple variables or complex logic → use a mutex.
//
// ============================================================================
// CORE FUNCTIONS
// ============================================================================
//
//   atomic.AddInt64(&val, 1)      — add 1 to val atomically (increment)
//   atomic.AddInt64(&val, -1)     — subtract 1 (decrement)
//   atomic.LoadInt64(&val)        — safely read val
//   atomic.StoreInt64(&val, 100)  — safely write val
//   atomic.CompareAndSwapInt64(&val, old, new)
//     — sets val to new ONLY if current value == old. Returns true if swapped.
//     — useful for lock-free updates where you only want to write if nothing
//       else changed the value since you last read it.
//
// WHY USE LoadInt64 INSTEAD OF JUST READING THE VARIABLE?
//   Reading a variable that another goroutine might be writing at the same
//   time is a data race — even if you're just reading.
//   atomic.LoadInt64() guarantees you get a consistent value.
//
// ============================================================================
// COMMON MISTAKES
// ============================================================================
//
// 1. Using & (address) correctly
//    All atomic functions take a POINTER to the variable, not the value.
//    atomic.AddInt64(&counter, 1)  ✓
//    atomic.AddInt64(counter, 1)   ✗ — won't compile
//
// 2. Using int instead of int64
//    The atomic package works with fixed-size types: int32, int64, uint32,
//    uint64, uintptr. It does NOT work with plain `int`. Use int64.
//
// 3. Mixing atomic and non-atomic access
//    If you use atomic.AddInt64() in one place but counter++ in another,
//    you still have a data race. ALL accesses must be atomic.
//    ✗  atomic.AddInt64(&c, 1) in one goroutine + c++ in another
//    ✓  atomic.AddInt64(&c, 1) everywhere
//
// 4. defer wg.Done() AFTER the loop, not at the top of the goroutine
//    The original code placed `defer wg.Done()` at the bottom of the
//    goroutine body (after the for loop). This works but is risky —
//    if the loop panics, defer still runs, but if you add an early return
//    later, you might miss it. Convention: always put defer wg.Done()
//    at the TOP of the goroutine, right after the opening brace.
//
// ============================================================================
// WHEN TO USE ATOMIC IN REAL SYSTEMS
// ============================================================================
//
//   - Request counters, hit counters, error counters in servers.
//   - Feature flags: a single bool/int that goroutines read frequently.
//   - Tracking active goroutine count or connection count.
//   - Any single shared number that only needs simple add/subtract/read.
//
// ============================================================================
// INTERVIEW POINTS
// ============================================================================
//
//   Q: What does "atomic" mean?
//   A: The operation is indivisible — it completes in one step with no
//      possibility of another goroutine interrupting it halfway.
//
//   Q: Why not just use a mutex for everything?
//   A: Atomic operations are faster for simple cases — no lock/unlock
//      overhead, no goroutine blocking. But they only work on single
//      variables with simple operations.
//
//   Q: Why does atomic.Load exist? Why not just read the variable?
//   A: Reading a variable while another goroutine writes it is a data race,
//      even if you're only reading. atomic.Load guarantees a safe read.
//
//   Q: What is CompareAndSwap used for?
//   A: Lock-free updates. Read a value, compute a new one, then only write
//      it if nobody else changed it in the meantime. If the swap fails,
//      retry. This is how many lock-free data structures work.
//
//   Q: Can you use atomic with a plain int?
//   A: No. The atomic package requires fixed-size types: int32 or int64.
//
// ============================================================================

type atomicCounter struct {
	count int64 // int64 required — atomic package does not support plain int
}

func (ac *atomicCounter) increment() {
	atomic.AddInt64(&ac.count, 1)
}

func (ac *atomicCounter) getVal() int64 {
	// LoadInt64 ensures a safe read — a plain `return ac.count` here
	// would be a data race if another goroutine is writing at the same time.
	return atomic.LoadInt64(&ac.count)
}

func runAtomicCounter() {
	var wg sync.WaitGroup
	numGoRoutines := 10
	c := &atomicCounter{}

	for range numGoRoutines {
		wg.Add(1)
		go func() {
			defer wg.Done() // always at the top, not the bottom
			for range 1000 {
				c.increment()
				// c.count++  ← data race: not atomic, not protected by mutex
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Final value: %d\n", c.getVal()) // always 10000
}

// ============================================================================
// COMPAREANDSWAP — practical example
// ============================================================================
//
// CAS (Compare and Swap) only updates the value if it still matches what
// you expect. If another goroutine changed it, the swap fails and you retry.
// This is how you do a lock-free update when the new value depends on the old.

func runCAS() {
	var val int64 = 0

	// Try to change val from 0 to 42 — only succeeds if val is still 0.
	swapped := atomic.CompareAndSwapInt64(&val, 0, 42)
	fmt.Printf("CAS swapped: %v, val: %d\n", swapped, atomic.LoadInt64(&val))

	// Try again — val is now 42, not 0, so this fails.
	swapped = atomic.CompareAndSwapInt64(&val, 0, 99)
	fmt.Printf("CAS swapped: %v, val: %d\n", swapped, atomic.LoadInt64(&val))
}

func AtomicCounter() {
	fmt.Println("=== ATOMIC COUNTER ===")
	runAtomicCounter()

	fmt.Println("\n=== COMPARE AND SWAP ===")
	runCAS()
}
