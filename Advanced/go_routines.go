package main

import (
	"fmt"
	"time"
)

// ====================================================
// GOROUTINES
// ====================================================
// WHAT    : Lightweight functions that run concurrently
//           on a separate goroutine, not the main thread
// KEY     : Non-blocking — do NOT pause program flow
// START   : go funcName() → launches goroutine immediately
// END     : goroutine exits when its function returns
// RISK    : Main goroutine exits → ALL goroutines are killed
// ====================================================

func Goroutines() {
	fmt.Println("Beginning Program")

	go sayHello() // go keyword → extracts function from main thread, runs concurrently

	fmt.Println("After SayHello()")
	// NOTE: "After SayHello()" prints BEFORE "Hello From Goroutine"
	// because sayHello() is non-blocking — main thread moves on immediately

	go printNumbers()
	go printLetters()
	// NOTE: Execution order of printNumbers vs printLetters is NOT guaranteed
	// The Go scheduler decides which runs first, when, and for how long

	var err error

	// err = go doWork() // INVALID — go returns no value to the caller
	go func() {
		// PATTERN: wrap in anonymous goroutine to capture return values
		// into variables shared with the outer scope
		err = doWork()
	}()

	// RACE CONDITION HERE ↓
	// Main goroutine reads  → err
	// Worker goroutine writes → err
	// Result depends on which runs first → undefined behaviour
	if err != nil {
		fmt.Println("Error Occured", err)
	} else {
		fmt.Println("Work Done Successfully")
	}
	// FIX: Use channels or sync.WaitGroup to wait for the goroutine
	// before reading shared data

	time.Sleep(2 * time.Second)
	// WHY Sleep() here: keeps main goroutine alive long enough
	// for background goroutines to finish
	// REAL FIX: sync.WaitGroup is the correct approach, not Sleep()
}

// ====================================================
// Sleep() BEHAVIOUR
// ====================================================
// - Pauses ONLY the current goroutine, not the whole program
// - Smaller duration → goroutine becomes runnable more often
//   → gets more scheduler opportunities
// - Does NOT control execution order — scheduler still decides
// ====================================================

func sayHello() {
	time.Sleep(1 * time.Second)
	fmt.Println("Hello From Goroutine")
}

func printNumbers() {
	for i := range 4 {
		fmt.Println(i, time.Now())
		time.Sleep(100 * time.Millisecond)
	}
}

func printLetters() {
	for _, i := range "abcd" {
		fmt.Println(string(i), time.Now())
		time.Sleep(200 * time.Millisecond)
		
	}
}

// ====================================================
// GOROUTINE + ERROR PATTERN
// ====================================================
// WRONG : err = go doWork()          → compiler error, no return from go
// RIGHT : go func() { err = doWork() }()
// BETTER: use channels to pass error back safely
//
//   ch := make(chan error)
//   go func() { ch <- doWork() }()
//   err := <-ch   ← blocks until goroutine sends
// ====================================================

func doWork() error {
	time.Sleep(1 * time.Second)
	fmt.Println("Checking it should be printed at last")
	return fmt.Errorf("An error occured")
}

// ====================================================
// CONCURRENCY vs PARALLELISM
// ====================================================
// Concurrency  : Multiple tasks make INDEPENDENT PROGRESS
//                (not necessarily at the same time)
// Parallelism  : Multiple tasks execute at the EXACT SAME TIME
//                (requires multiple CPU cores)
//
// Goroutines   : facilitate concurrency
// Go runtime   : schedules goroutines across available CPUs
//                → enables parallelism automatically
// ====================================================

// ====================================================
// M:N SCHEDULING MODEL
// ====================================================
// M goroutines mapped onto N OS threads mapped onto CPU cores
//
//   Many Goroutines (M)
//          ↓  (Go scheduler)
//   Few OS Threads (N)
//          ↓
//       CPU Cores
//
// WHY : Creating OS threads is expensive
//       Goroutines are cheap (~2KB stack, grows as needed)
//       Go scheduler multiplexes them efficiently
// ====================================================

/*
========================================
QUICK REVISION
========================================

CONCEPT             KEY POINT
-----------         -------------------------------------------
go keyword          Launches function as goroutine (non-blocking)
Execution order     NOT guaranteed — scheduler decides
main exits          All goroutines are killed immediately
Sleep()             Pauses current goroutine only, not whole program
Race condition      Multiple goroutines access shared data → result
                    depends on timing → undefined behaviour
Concurrency         Tasks progress independently (not same time)
Parallelism         Tasks run at literally the same time (multi-CPU)
M:N Scheduling      M goroutines on N threads on CPU cores

----------------------------------------------------
PATTERNS

Launch goroutine        go funcName()
Capture return value    go func() { err = doWork() }()
Safe error return       ch := make(chan error)
                        go func() { ch <- doWork() }()
                        err := <-ch
Wait for goroutines     sync.WaitGroup (preferred over Sleep)

----------------------------------------------------
INTERVIEW QUESTIONS

Q1. What is a goroutine?
Q2. Is goroutine execution order guaranteed?
Q3. What happens when the main goroutine exits?
Q4. What does Sleep() pause — goroutine or whole program?
Q5. What is a race condition? Give an example from this file.
Q6. How do you return an error from a goroutine?
Q7. Concurrency vs Parallelism — what's the difference?
Q8. What is M:N scheduling?
Q9. Why are goroutines cheaper than OS threads?

========================================
*/