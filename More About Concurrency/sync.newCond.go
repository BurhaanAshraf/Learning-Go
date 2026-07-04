package main

import (
	"fmt"
	"sync"
	"time"
)

const bufferSize = 5

/*
================================================================
CORE CONCEPT — sync.Cond (Condition Variables)
================================================================

WHY IT EXISTS:
  A Mutex only answers "can I access this safely right now?" —
  it says nothing about WHEN the right moment to proceed is.
  sync.Cond adds that missing piece: a goroutine can sleep until
  some condition becomes true, without busy-waiting (looping and
  repeatedly checking, which burns CPU for nothing).

  Mutex              -> protects shared data
  Condition Variable -> coordinates WHEN goroutines should proceed

  Always paired with a Mutex — Cond has no data protection of its
  own; it relies on the Locker you give it via sync.NewCond().

METHODS:
  Wait()      Sleep until signaled. Internally: Unlock -> sleep ->
              (on wake) re-Lock before returning to your code.
  Signal()    Wake ONE waiting goroutine (if any are waiting).
  Broadcast() Wake ALL waiting goroutines.

USE CASES: producer-consumer, task queues, resource pools —
  anywhere a goroutine needs to wait for "something to change"
  rather than just "access being free".
================================================================
*/

type Buffer struct {
	items []int      // shared queue
	mutex sync.Mutex // protects items
	cond  *sync.Cond // coordinates producers/consumers waiting on items
}

// newBuffer: sync.NewCond needs a Locker (anything with Lock/Unlock).
// sync.Mutex satisfies that interface, so we pass &b.mutex. Cond
// can't be set in the struct literal above because it needs a
// pointer to b.mutex, which doesn't exist until b itself does.
func newBuffer(size int) *Buffer {
	b := &Buffer{items: make([]int, 0, size)}
	b.cond = sync.NewCond(&b.mutex)
	return b
}

// produce: blocks (via Wait) while the buffer is full, adds an
// item once space frees up, then wakes one waiting consumer.
func (b *Buffer) produce(item int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// MUST be a `for` loop, never `if`. Two reasons:
	// 1. Even after being woken, the condition might no longer
	//    hold (e.g. another producer snuck in and filled it again
	//    before this goroutine got the mutex back).
	// 2. Spurious wakeups: Go's Cond.Wait() can return even
	//    without anyone calling Signal()/Broadcast(). Re-checking
	//    the condition is the only safe pattern.
	for len(b.items) == bufferSize {
		b.cond.Wait() // atomically: unlock -> sleep -> re-lock on wake
	}

	b.items = append(b.items, item)
	fmt.Println("Produced:", item)

	b.cond.Signal() // wake one waiting consumer (buffer now has data)
}

// consume: blocks while the buffer is empty, removes the oldest
// item once available, then wakes one waiting producer.
func (b *Buffer) consume() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for len(b.items) == 0 {
		b.cond.Wait()
	}

	item := b.items[0]
	b.items = b.items[1:]
	fmt.Println("Consumed:", item)

	b.cond.Signal() // wake one waiting producer (buffer now has space)

	return item
}

func producer(b *Buffer, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range 10 {
		b.produce(i + 100)
		time.Sleep(100 * time.Millisecond) // demo pacing only
	}
}

func consumer(b *Buffer, wg *sync.WaitGroup) {
	defer wg.Done()
	for range 10 {
		b.consume()
		time.Sleep(100 * time.Millisecond) // demo pacing only
	}
}

func newCond() {
	// Producer -> Buffer -> Consumer
	// Producer waits (Wait()) when buffer is full.
	// Consumer waits (Wait()) when buffer is empty.
	// sync.Cond coordinates both sides without busy-waiting.

	buffer := newBuffer(bufferSize)
	var wg sync.WaitGroup

	wg.Add(2)
	go producer(buffer, &wg)
	go consumer(buffer, &wg)
	wg.Wait()

	fmt.Println("Main Gracefully Executed...")
}

/*
================================================================
WAIT() MECHANISM (same shape for both producer and consumer)
================================================================
Condition met (space free / item available)?
  ├── Yes -> proceed (produce or consume)
  └── No  -> Wait()
               ↓
         Unlock mutex
               ↓
         Sleep goroutine
               ↓
      Signal()/Broadcast() (or spurious wakeup)
               ↓
         Re-lock mutex
               ↓
       Loop back -> re-check condition
================================================================
QUICK REVISION
================================================================
sync.Cond    -> coordinates goroutines waiting on a condition.
Wait()       -> unlock -> sleep -> re-lock -> return.
Signal()     -> wakes ONE waiting goroutine.
Broadcast()  -> wakes ALL waiting goroutines.

Always:  for condition { cond.Wait() }
Never:   if condition { cond.Wait() }
Why: condition may change again before you reacquire the mutex,
     AND Wait() can return from a spurious wakeup with no Signal
     having been sent at all.

Mutex               -> protects shared DATA
Condition Variable  -> coordinates EXECUTION ORDER
================================================================
INTERVIEW Q&A
================================================================
Q: Why use sync.Cond instead of just a Mutex + busy-loop?
A: Busy-looping wastes CPU checking a condition repeatedly.
   Cond lets the goroutine sleep until actually signaled.

Q: Difference between Wait(), Signal(), and Broadcast()?
A: Wait() sleeps until signaled. Signal() wakes one waiter.
   Broadcast() wakes all waiters — use when multiple goroutines
   could all validly proceed (vs Signal() for just one).

Q: Why is Wait() always inside a for loop instead of an if?
A: The condition can become false again before this goroutine
   reacquires the mutex, and spurious wakeups can trigger Wait()
   to return with no signal sent at all. Re-check, don't assume.

Q: What happens internally when Wait() is called?
A: Atomically unlocks the mutex and puts the goroutine to sleep;
   on wake, it reacquires the mutex before returning to the caller.

Q: Can sync.Cond be used without a Mutex?
A: No — it requires a Locker (Mutex or RWMutex) passed to
   sync.NewCond(); it has no data-protection mechanism of its own.
================================================================
*/
