package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ACTOR PATTERN: one goroutine exclusively owns and mutates state (count).
// Others communicate through a channel. No mutex needed — no contention exists.
//
// Go proverb: "Don't share memory to communicate. Communicate to share memory."
//
// WHY THIS MATTERS — without this pattern, two goroutines doing count++ simultaneously:
//   A reads 5, B reads 5 → A writes 6, B writes 6 → expected 7, got 6. Silent bug.

type statefulWorker struct {
	count int      // mutable state; only the worker goroutine may read/write it
	ch    chan int // inbox — everyone else sends values through here
}

// Start launches the goroutine that owns state. Returns immediately.
// The goroutine runs independently and continuously processes messages from ch.
func (st *statefulWorker) Start(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {

			// select waits until one of its channel operations can proceed.
			// With a single case, it behaves like a blocking receive.
			// Add more cases later (e.g. quit signal) without restructuring.
			select {
			case val := <-st.ch:
				st.count += val // safe — only this goroutine ever writes count
				fmt.Println("Current Count:", st.count)

			case <-ctx.Done():
				return
			}
		}
	}() // go func(){ ... }() — anonymous func, launched immediately as a goroutine
}

// Send blocks until the worker receives the value.
// Unbuffered channel forces a synchronous handoff — sender can't outrun the worker.
func (st *statefulWorker) Send(value int) {
	// Senders never modify count directly.
	// They can only request changes by sending messages.
	st.ch <- value
}

func stateful_goroutines() {

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	st := &statefulWorker{
		ch: make(chan int), // unbuffered: sender and receiver must meet at the same time
		// count zero-initialized by Go automatically
	}

	wg.Add(1)
	st.Start(ctx, &wg)

	// range 5 → i = 0,1,2,3,4  (Go 1.22+). Final count = 0+1+2+3+4 = 10
	for i := range 5 {
		st.Send(i)
		time.Sleep(500 * time.Millisecond)

	}

	cancel()
	wg.Wait()
	fmt.Println("Closing Channel Gracefully...")
	fmt.Println("Execution Completed")
}

// - One goroutine owns the state.
// - Other goroutines communicate with it through channels.
// - State is never shared directly.

// In this example:

// Main Goroutine
//     |
//     v
//  Channel
//     |
//     v
// Worker Goroutine (owns count)

// Key Rules:

// 1. Only the worker goroutine can modify count.
// 2. Main sends requests through the channel.
// 3. No shared mutable state → no mutex needed.
// 4. Messages travel through the channel; state stays with the worker.
// 5. Shutdown is coordinated using context and WaitGroup.

// Flow:

// Send(value) -> Channel -> Worker -> Update State

// Mental Model:

// The worker is a manager holding the state.
// The channel is the manager's inbox.
// Everyone else submits requests instead of touching the state directly.
