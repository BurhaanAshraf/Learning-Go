package main

import (
	"fmt"
	"sync"
	"time"
)

/*
================================================================
CORE CONCEPT — sync.RWMutex
================================================================

WHY IT EXISTS:
  A plain sync.Mutex only allows ONE goroutine in at a time —
  even for reads. But reads don't modify data, so two readers
  running simultaneously can never corrupt anything. A plain
  Mutex forces them to queue up anyway, wasting concurrency.

  RWMutex fixes this by splitting locking into two modes:
    RLock() / RUnlock()  -> shared lock, for reads
    Lock()   / Unlock()  -> exclusive lock, for writes

COMPATIBILITY TABLE:
              RLock()      Lock()
  RLock()       ✓            ✗
  Lock()        ✗            ✗

  - Many readers            -> allowed together
  - One writer               -> allowed alone
  - Reader + Writer together -> never allowed

WHY THIS IS SAFE:
  Lock() waits until every existing RLock() is released, then
  blocks all new RLock()/Lock() calls until it's done. This
  guarantees a writer never sees readers mid-read, and readers
  never see a half-written value.

WHEN TO USE IT:
  Use RWMutex when reads >> writes (e.g. config, cache, shared
  maps read constantly but updated rarely). If writes are
  frequent, the RLock/Lock bookkeeping overhead isn't worth it —
  a plain Mutex is simpler and often just as fast.
================================================================
*/

var (
	rwmu    sync.RWMutex
	counter int // shared state — always access through rwmu
)

// readCounter: RLock() lets this run CONCURRENTLY with other
// readers — safe because reading never mutates shared state.
func readCounter(wg *sync.WaitGroup) {
	defer wg.Done()

	rwmu.RLock()
	fmt.Println("Read Counter:", counter)
	rwmu.RUnlock()
}

// writeCounter: Lock() is exclusive — blocks until all current
// readers finish, then blocks any new readers/writers until
// this Unlock()s. Prevents readers from seeing a half-written value.
func writeCounter(wg *sync.WaitGroup, value int) {
	defer wg.Done()

	rwmu.Lock()
	counter = value
	fmt.Println("Updated Counter:", value)
	rwmu.Unlock()
}

func RWMutex() {
	var wg sync.WaitGroup
	fmt.Println("Starting Main...")

	// 5 readers, all holding RLock() at once — this is the
	// scenario RWMutex is optimized for, and where a plain
	// Mutex would've forced needless serialization.
	for range 5 {
		wg.Add(1)
		go readCounter(&wg)
	}

	time.Sleep(1 * time.Second) // demo-only: lets readers finish before the writer starts

	// Single writer — Lock() here waits for all 5 RLock()s
	// above to release before it can proceed.
	wg.Add(1)
	go writeCounter(&wg, 18)

	wg.Wait()
	fmt.Println("Main Gracefully Executed...")
}

/*
================================================================
1-MINUTE REVISION
================================================================
Mutex     -> one goroutine at a time, period (read or write).
RWMutex   -> many readers together, OR one writer alone, never both.

RLock()  = shared read lock   (concurrent readers OK)
Lock()   = exclusive write lock (blocks everyone else)

Compatibility:
  Reader + Reader -> ✓
  Reader + Writer -> ✗
  Writer + Writer -> ✗

Use when: reads >> writes (cache, config, shared maps).
Skip when: writes are frequent -> plain Mutex is simpler.
================================================================
INTERVIEW Q&A
================================================================
Q: Mutex vs RWMutex?
A: Mutex = one goroutine at a time. RWMutex = many readers,
   but only one writer, never mixed.

Q: When to prefer RWMutex?
A: When reads happen far more often than writes.

Q: Can multiple goroutines hold RLock() at once?
A: Yes.

Q: Can a writer get Lock() while readers hold RLock()?
A: No — it waits until every reader releases RUnlock().

Q: Can a reader get RLock() while a writer holds Lock()?
A: No — it waits until the writer calls Unlock().
================================================================
*/
