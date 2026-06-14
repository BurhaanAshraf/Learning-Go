package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// sync.WaitGroup — REVISION NOTES
// ============================================================================
//
// WHAT
//   A WaitGroup is a counter. You tell it how many goroutines to wait for,
//   each goroutine signals when it's done, and your main code blocks until
//   the counter hits zero.
//
// WHY
//   go worker() starts a goroutine but doesn't wait for it. If main()
//   returns before the goroutine finishes, the goroutine is killed mid-work.
//   WaitGroup is how you wait for a group of goroutines to finish.
//
// THE 3 METHODS
//   wg.Add(n)   — increase the counter by n. Call this BEFORE starting
//                 the goroutines.
//   wg.Done()   — decrease the counter by 1. Call this inside the
//                 goroutine when it's done. Usually written as
//                 `defer wg.Done()` at the top of the function.
//   wg.Wait()   — blocks until the counter is 0.
//
// MENTAL MODEL
//   Add(3) → counter = 3
//   each Done() → counter -= 1
//   Wait() unblocks when counter == 0
//
// ALWAYS PASS BY POINTER
//   func Worker(id int, wg *sync.WaitGroup)
//   WaitGroup must be passed as *sync.WaitGroup. If you pass it by value,
//   each goroutine gets its own copy of the counter — Wait() will never
//   see the real count and either returns immediately or never unblocks.
//
// ============================================================================
// COMMON MISTAKES
// ============================================================================
//
// 1. Forgetting wg.Done()
//    Counter never reaches 0. wg.Wait() blocks forever (deadlock).
//
// 2. Calling wg.Add() AFTER starting the goroutine
//    Race condition — Wait() might run before Add(), so it sees counter=0
//    and returns immediately, before the goroutine even starts.
//    Always Add() before `go func...`.
//
// 3. Passing WaitGroup by value, not by pointer
//    Each goroutine modifies its own copy. Wait() on the original never
//    sees Done() calls.
//
// 4. Calling wg.Done() more times than wg.Add()
//    Counter goes negative → panic: "negative WaitGroup counter".
//
// ============================================================================

// Worker is a simple goroutine that does some "work" and signals when done.
func Worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // runs when Worker returns, however it returns

	time.Sleep(1 * time.Second)
	fmt.Printf("Worker %d starting...\n", id)

	time.Sleep(2 * time.Second) // simulate work
	fmt.Printf("Worker %d has finished the task...\n", id)
}

func runWorker() {
	var wg sync.WaitGroup
	numWorkers := 3

	wg.Add(numWorkers) // Add BEFORE launching goroutines

	for i := range numWorkers {
		go Worker(i+1, &wg)
	}

	wg.Wait() // blocks until all 3 call Done()
	fmt.Println("All workers finished...")
}

// ============================================================================
// WAITGROUP + CHANNELS — closing `results` safely
// ============================================================================
//
// In a worker pool, multiple workers write to `results`. Only ONE of them
// should close it, and only after ALL of them are done writing.
//
// PATTERN:
//   wg.Add(numWorkers)
//   start workers (each does `defer wg.Done()`)
//   go func() {
//       wg.Wait()      // wait for all workers to finish writing
//       close(results) // now it's safe to close
//   }()
//   for r := range results { ... }  // exits once results is closed+empty
//
// WHY THE close() IS IN ITS OWN GOROUTINE
//   wg.Wait() blocks. If we called it in main before the `for range results`
//   loop, we'd deadlock: workers can't finish sending to `results` if main
//   isn't reading it yet (once the buffer fills up), but main is stuck on
//   Wait() instead of reading. Running Wait()+close() in a separate
//   goroutine lets main read results WHILE workers are still finishing.
//
// ============================================================================

func withChannels(id int, tasks <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d is starting...\n", id)
	time.Sleep(500 * time.Millisecond) // simulate startup cost

	for task := range tasks {
		results <- task * 4
	}
	fmt.Printf("Worker %d has finished...\n", id)
}

func runWithChannels() {
	numJobs := 6
	numWorkers := 3

	tasks := make(chan int, numJobs)
	results := make(chan int)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// start workers FIRST — they'll wait on `range tasks` until jobs arrive
	for i := range numWorkers {
		go withChannels(i+1, tasks, results, &wg)
	}

	// send jobs
	for i := range numJobs {
		tasks <- i + 1
	}
	close(tasks) // tells workers: no more jobs, exit range loop after draining

	// close results only after all workers are done writing to it
	go func() {
		wg.Wait()
		close(results)
	}()

	for range results {
		fmt.Println("Values Received:", <-results)
	}
	fmt.Println("Execution Successful...")
}

// ============================================================================
// WAITGROUP WITH METHODS — one goroutine per task, not a fixed pool
// ============================================================================
//
// Earlier examples had a FIXED pool of workers pulling from a shared queue.
// Here, every task gets its OWN goroutine — useful when each task is
// independent and you don't need to limit concurrency.
//
// wg.Add(1) is called ONCE PER TASK, right before starting that task's
// goroutine — not once with the total count up front. Same end result,
// but this style is common when the number of tasks isn't known ahead
// of time (e.g. reading from a stream).
// ============================================================================

type WorkerDetails struct {
	ID   int
	TASK string
}

// PerformTask is a method on WorkerDetails — same role as the Worker
// function above, just attached to a struct so it has access to ID/TASK.
func (w *WorkerDetails) PerformTask(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("WorkerID %d started working on task %s...\n", w.ID, w.TASK)
	time.Sleep(2 * time.Second)
	fmt.Printf("WorkerID %d finished task %s...\n", w.ID, w.TASK)
}

func construction() {
	var wg sync.WaitGroup
	tasks := []string{"digging", "laying bricks", "painting", "furniture", "putting glasses"}

	for i, task := range tasks {
		worker := WorkerDetails{ID: i + 1, TASK: task}
		wg.Add(1) // one Add per task, right before launching it
		go worker.PerformTask(&wg)
	}

	wg.Wait()
	fmt.Println("Construction is Finished...")
}

// ============================================================================
// INTERVIEW POINTS
// ============================================================================
//
//   Q: What does WaitGroup do?
//   A: Lets you wait for a group of goroutines to finish before continuing.
//
//   Q: Why pass *sync.WaitGroup, not sync.WaitGroup?
//   A: Passing by value gives each goroutine its own copy of the counter —
//      the original Wait() never sees the Done() calls.
//
//   Q: Why call Add() before `go func...`, not inside it?
//   A: If Wait() runs before Add(), it sees counter=0 and returns
//      immediately — possibly before any work has started.
//
//   Q: How do you safely close a channel that multiple goroutines write to?
//   A: Use a WaitGroup. Have a separate goroutine call wg.Wait() then
//      close(channel), once all writers are confirmed done.
//
//   Q: wg.Add(n) once vs wg.Add(1) per task — when would you use each?
//   A: Add(n) up front when you know the count in advance (fixed pool).
//      Add(1) per task when tasks are created dynamically / one at a time.
//
// ============================================================================

func WaitGroups() {
	fmt.Println("=== BASIC WAITGROUP ===")
	runWorker()

	fmt.Println("\n=== WAITGROUP + CHANNELS ===")
	runWithChannels()

	fmt.Println("\n=== WAITGROUP PER TASK (METHOD) ===")
	construction()
}
