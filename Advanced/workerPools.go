package main

import (
	"fmt"
	"time"
)

// ============================================================================
// WORKER POOLS — REVISION NOTES
// ============================================================================
//
// WHAT
//   A fixed number of goroutines (workers) read jobs from a channel and
//   process them, one job at a time per worker.
//
// WHY
//   100 jobs do not need 100 goroutines.
//   A worker pool uses a fixed number of goroutines to handle many jobs.
//   This caps how much runs at once — limits memory, CPU, and load on
//   things like databases or APIs.
//
// HOW IT WORKS
//   - tasks channel   = jobs go in  (producer writes, workers read)
//   - results channel = results come out (workers write, producer reads)
//   - Each worker runs: for job := range tasks { do work }
//   - close(tasks) tells workers "no more jobs" — their loops exit.
//
// ============================================================================
// SKELETON (memorize this — 5 steps)
// ============================================================================
//
//   1. Make two channels: one for jobs in, one for results out.
//        tasks := make(chan int, 10)
//        results := make(chan int, 10)
//
//   2. Start the workers. They wait for jobs immediately.
//        for i := 1; i <= 3; i++ {
//            go worker(i, tasks, results)
//        }
//
//   3. Send the jobs.
//        for i := 1; i <= 10; i++ {
//            tasks <- i
//        }
//
//   4. Close tasks — tells workers "no more jobs."
//        close(tasks)
//
//   5. Read the results (we know there are 10).
//        for i := 0; i < 10; i++ {
//            fmt.Println(<-results)
//        }
//
// ============================================================================
// GOROUTINES — QUICK RECAP
// ============================================================================
//
// WHAT
//   A goroutine is a lightweight function that runs at the same time as
//   other code. Start one by writing `go` before a function call.
//
//   go doSomething()   // runs concurrently, doesn't block the caller
//
// WHY
//   Lets your program do multiple things at once — e.g. handle many
//   requests, or process many jobs in parallel.
//
// GOTCHA
//   If main() returns before a goroutine finishes, that goroutine is
//   killed mid-work — its output may never appear. Channels (below) are
//   how we wait for goroutines properly.
//
// ============================================================================
// CHANNELS — QUICK RECAP
// ============================================================================
//
// WHAT
//   A channel is a typed pipe. One goroutine sends a value in, another
//   goroutine receives it out.
//
//     ch := make(chan int)  // unbuffered channel of ints
//     ch <- 5                // send 5 into the channel
//     x := <-ch              // receive a value from the channel
//
// UNBUFFERED vs BUFFERED
//   - Unbuffered (make(chan int)):
//       A send blocks until someone is ready to receive right now.
//   - Buffered (make(chan int, 5)):
//       Can hold up to 5 values before a send blocks.
//       Useful for letting a producer get ahead of consumers.
//
//   In real systems, buffer size is usually smaller than total jobs.
//   When the buffer is full, sends block — this slows the producer down
//   to match the workers. This is called backpressure, and it's a good thing.
//
// ============================================================================
// close() + range — HOW WORKERS KNOW WHEN TO STOP
// ============================================================================
//
//   for job := range tasks { ... }
//
//   - Waits when tasks is empty but still open.
//   - Runs the loop body when a job arrives.
//   - Exits with no error once tasks is closed AND empty.
//
// RULES
//   - Only the sender closes a channel. Only once.
//   - Send on a closed channel → panic.
//   - Close a channel twice → panic.
//   - Read from a closed+empty channel → zero value, instantly, no block.
//
// WHY `results` IS NOT CLOSED HERE
//   We know exactly how many results to expect, so we just read that many.
//   No need to close it.
//
// ============================================================================
// CHANNEL DIRECTION IN FUNCTION SIGNATURES
// ============================================================================
//
//   func worker(id int, tasks <-chan int, results chan<- int)
//
//   <-chan int   = this function can only READ from tasks
//   chan<- int   = this function can only WRITE to results
//
// WHY
//   It documents intent and the compiler enforces it — a worker can't
//   accidentally send on `tasks` or read from `results`. Bugs caught at
//   compile time instead of runtime.
//
// ============================================================================
// COMMON MISTAKES
// ============================================================================
//
// 1. Forgetting close(tasks)
//    Workers wait forever for the next job. They never exit. Goroutines leak.
//
// 2. A worker closes tasks
//    Only the producer should close it. If anything else closes or sends
//    on it after, the program panics.
//
// 3. Reading more results than were sent
//    Extra <-results blocks forever → deadlock.
//
// 4. Assuming results come back in the same order jobs were sent
//    With multiple workers, faster jobs finish first.
//    If order matters, put an ID in the job and send that ID back with
//    the result (see ticketRequest below).
//
// 5. Using a loop variable directly inside a goroutine
//    for i := 1; i <= 3; i++ {
//        go func() { fmt.Println(i) }()   // BAD: may print wrong/same i
//    }
//    Fix: pass it in as a parameter.
//        go func(id int) { fmt.Println(id) }(i)   // GOOD
//
// ============================================================================
// BEST PRACTICES
// ============================================================================
//
//   - Use <-chan T / chan<- R in function params (see "Channel Direction" above).
//
//   - Worker count:
//       CPU-heavy work          → about as many workers as CPU cores.
//       I/O work (DB, API calls) → can use more workers than cores,
//       since workers spend most time waiting, not using CPU.
//
//   - Always have a clear "no more work" signal (close(tasks)) and a clear
//     "how many results to expect" so the program doesn't hang forever.
//
// ============================================================================
// REAL-WORLD USE
// ============================================================================
//
//   - Processing a batch of DB rows or file lines with limited concurrency.
//   - Background job processing in a server (e.g. sending emails).
//   - Pipelines: one stage's output feeds the next stage's worker pool.
//   - Calling an external API with a max number of concurrent requests.
//
// ============================================================================
// INTERVIEW POINTS
// ============================================================================
//
//   Q: Why not just `go` every job?
//   A: Too many goroutines at once. A pool caps concurrency.
//
//   Q: How do workers know to stop?
//   A: close(tasks). Their range loops exit once it's closed and empty.
//
//   Q: Who closes a channel?
//   A: Only the sender, only once.
//
//   Q: Does result order match job order?
//   A: No — results come back in completion order. Use an ID to match them.
//
//   Q: How do you know all workers have actually exited (not just
//      that you got all results)?
//   A: Getting all results just means the work is done, not that the
//      goroutines returned. To wait for that, use sync.WaitGroup —
//      each worker calls wg.Done(), main calls wg.Wait().
//
//   Q: Unbuffered vs buffered channel — when would you use each?
//   A: Unbuffered when you want sender and receiver to sync up exactly
//      (handshake). Buffered when you want the producer to not wait for
//      every single send — useful for worker pools and queues.
//
// ============================================================================

// worker reads ints from tasks, doubles them, sends to results.
// Exits when tasks is closed and empty.
func worker(id int, tasks <-chan int, results chan<- int) {
	for task := range tasks {
		fmt.Printf("[Worker %d] Processing Task %d\n", id, task)
		time.Sleep(500 * time.Millisecond) // simulate work (e.g. API call)
		results <- task * 2
	}
	fmt.Printf("[Worker %d] No more tasks. Exiting.\n", id)
}

func runBasicWorkerPool() {
	const (
		numberOfWorkers = 3
		numberOfJobs    = 10
	)

	tasks := make(chan int, numberOfJobs)
	results := make(chan int, numberOfJobs)

	for i := 1; i <= numberOfWorkers; i++ {
		go worker(i, tasks, results)
	}

	for i := 1; i <= numberOfJobs; i++ {
		tasks <- i
	}
	close(tasks) // no more jobs coming

	for i := 0; i < numberOfJobs; i++ {
		result := <-results
		fmt.Println("Result received:", result)
	}
}

// ============================================================================
// EXAMPLE 2: same pattern, struct instead of int.
// Shows how to track which job a result belongs to (see mistake #4).
// ============================================================================

type ticketRequest struct {
	id         int
	numTickets int
	cost       int
}

func ticketProcessing(workerID int, requests <-chan ticketRequest, results chan<- int) {
	for req := range requests {
		fmt.Printf("[TicketWorker %d] Processing %d ticket(s) for personID %d (Cost: $%d)\n",
			workerID, req.numTickets, req.id, req.cost)
		time.Sleep(400 * time.Millisecond)
		results <- req.id // send back the ID so we know which request finished

	}

}

func runTicketSystem(done chan<- struct{}) {
	const (
		numRequests     = 5
		price           = 5
		numberOfWorkers = 3
	)
	ticketRequests := make(chan ticketRequest, numRequests)
	ticketResults := make(chan int, numRequests)

	for i := 1; i <= numberOfWorkers; i++ {
		go func(workerID int) {
			ticketProcessing(workerID, ticketRequests, ticketResults)
		}(i)
	}

	for i := 1; i <= numRequests; i++ {
		ticketRequests <- ticketRequest{
			id:         i,
			numTickets: i * 2,
			cost:       i * price,
		}
	}
	close(ticketRequests)

	for i := 0; i < numRequests; i++ {
		completedID := <-ticketResults
		fmt.Printf("Ticket for PersonID %d processed completely\n", completedID)
	}

	done <- struct{}{}
}

func workerPools() {
	done := make(chan struct{})
	fmt.Println("=== RUNNING BASIC WORKER POOL ===")
	runBasicWorkerPool()

	fmt.Println("\n=== RUNNING TICKET SYSTEM ===")
	go runTicketSystem(done)

	<-done

}
