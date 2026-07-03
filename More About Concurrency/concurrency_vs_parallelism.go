package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

/*
================================================================
CORE CONCEPT — Concurrency vs Parallelism
================================================================

CONCURRENCY = dealing with multiple tasks at once (a DESIGN property).
  Tasks are structured to make progress independently, but they may
  or may not run at the literal same instant. Think: one chef
  switching between 3 dishes — never doing two at once, but all
  3 are "in progress" concurrently.

PARALLELISM = executing multiple tasks at the literal same instant
  (an EXECUTION property). Requires multiple CPU cores. Think:
  3 chefs, each cooking one dish, simultaneously.

KEY INSIGHT: concurrency is about STRUCTURE, parallelism is about
  HARDWARE. You can have concurrency with zero parallelism
  (1 core, many goroutines, scheduler time-slices between them).
  You cannot have parallelism without concurrency (you need
  independent tasks to run in parallel in the first place).

  Rob Pike's famous line: "Concurrency is about dealing with lots
  of things at once. Parallelism is about doing lots of things
  at once." — concurrency is the more general concept.

HOW GO ACHIEVES EACH:
  - Concurrency comes for free from goroutines + the Go scheduler
    (M:N model — many goroutines mapped onto few OS threads).
  - Parallelism additionally requires GOMAXPROCS > 1 AND actual
    multiple physical cores. GOMAXPROCS is the ceiling on how many
    goroutines can run Go code AT THE SAME INSTANT.
  - Default GOMAXPROCS = runtime.NumCPU() since Go 1.5 — you rarely
    need to set it manually.

WHY THIS DEMO IS SPLIT INTO TWO PARTS:
  Part 1 (printNumbers/printLetters) shows CONCURRENCY without
  needing parallelism — output interleaves because of Sleep()
  yielding, not because two cores are racing. Even on 1 core,
  you'd see the same interleaved output.

  Part 2 (HeavyTask) shows real PARALLELISM — a CPU-bound loop
  with no yielding, spread across GOMAXPROCS cores, actually
  executing simultaneously. This is where multiple cores matter.
================================================================
*/

// printNumbers/printLetters demonstrate CONCURRENCY:
// two independent tasks making interleaved progress. The
// interleaving happens because Sleep() yields control back to
// the scheduler — NOT because they're running in parallel.
// (Proof: this looks identical whether GOMAXPROCS is 1 or 12.)
func printNumbers(done chan struct{}) {
	for i := range 5 {
		fmt.Println(time.Now())
		fmt.Println(i)
		time.Sleep(500 * time.Millisecond) // yield point -> scheduler can switch to printLetters
	}
	close(done) // signal: "I'm finished"
}

func printLetters(done chan struct{}) {
	for _, i := range "ABCDE" {
		fmt.Println(time.Now())
		fmt.Println(string(i))
		time.Sleep(500 * time.Millisecond)
	}
	close(done)
}

// HeavyTask demonstrates PARALLELISM: a CPU-bound busy loop with
// no natural yield points. With GOMAXPROCS >= number of cores,
// multiple HeavyTask calls genuinely execute at the same instant
// on different cores — this is the part that actually benefits
// from having more CPUs.
func HeavyTask(id int, wg *sync.WaitGroup) {
	defer wg.Done() // deferred right after Add() so it runs even on panic

	fmt.Printf("Task %d is starting...\n", id)
	for range 100_000_000 {
	} // empty loop — fine for demo, not a real benchmark (compiler may optimize)
	fmt.Println(time.Now())
	fmt.Printf("Task %d is finished...\n", id)
}

func CVP() {
	// --- Part 1: CONCURRENCY demo — channels as completion signals ---
	// Using chan struct{} instead of Sleep(3*time.Second): we wait
	// for the EXACT moment each goroutine finishes, not a guessed
	// duration. struct{} carries no data — we only care THAT it's
	// done, not what value it produced.
	numbersDone := make(chan struct{})
	lettersDone := make(chan struct{})

	go printNumbers(numbersDone)
	go printLetters(lettersDone)

	<-numbersDone // blocks until printNumbers closes its channel
	<-lettersDone // blocks until printLetters closes its channel

	// --- Part 2: PARALLELISM demo — WaitGroup + multiple cores ---
	numThreads := 12
	runtime.GOMAXPROCS(numThreads) // parallelism ceiling: max goroutines running Go code at once

	var wg sync.WaitGroup
	for i := range numThreads {
		wg.Add(1)            // must happen before `go`, in this goroutine (not inside child)
		go HeavyTask(i, &wg) // &wg — WaitGroup must never be copied
	}
	wg.Wait() // blocks until all Done() calls bring counter to 0

	fmt.Println("Main Gracefully Executed...")
}

/*
================================================================
QUICK SELF-TEST TO CONFIRM THE CONCEPT (try it!):
  Change runtime.GOMAXPROCS(numThreads) to runtime.GOMAXPROCS(1)
  and re-run:
  - Part 1 (numbers/letters) output pattern barely changes —
    proves it never depended on parallelism.
  - Part 2 (HeavyTask) becomes visibly SLOWER / more serialized —
    all 12 tasks now compete for 1 core instead of running
    simultaneously across up to 12 cores. That slowdown IS
    parallelism (or its absence) made visible.
================================================================
*/
