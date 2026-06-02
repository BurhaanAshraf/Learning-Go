package main

import (
	"fmt"
	"time"
)

// ============================================================
// CHANNELS IN GO — Complete Notes & Runnable Code
// ============================================================
//
// WHAT IS A CHANNEL?
//
// A channel is a typed conduit for sending values between goroutines.
// It is Go's idiomatic way to share data safely across concurrent functions.
//
// Instead of:
//   Goroutine A ──► Shared Variable ◄── Goroutine B   (race condition risk)
//
// Go encourages:
//   Goroutine A ──► Channel ──► Goroutine B            (safe by design)
//
// Go Philosophy:
//   "Do not communicate by sharing memory;
//    share memory by communicating."
//
// ============================================================
// SYNTAX
// ============================================================
//
//   make(chan T)       → create an unbuffered channel of type T
//   ch <- value        → send value into channel  (BLOCKS until receiver is ready)
//   value := <-ch      → receive value from channel (BLOCKS until sender sends)
//
// ============================================================
// UNBUFFERED CHANNEL RULE
// ============================================================
//
//   Sender waits for receiver.
//   Receiver waits for sender.
//   The handoff happens only when BOTH sides are ready.
//   There is no internal queue — it is a direct synchronous handoff.
//
// ============================================================
// DEADLOCK
// ============================================================
//
//   ch := make(chan string)
//   ch <- "hello"
// fatal error: all goroutines are asleep - deadlock!
//
//   Cause: no receiver goroutine exists, so the sender blocks forever.
//   When ALL goroutines are blocked, the Go runtime detects it and panics.
//   Rule: every send must have a matching receive, and vice versa.
//
// ============================================================
// PRODUCER / CONSUMER PATTERN
// ============================================================
//
//   Producer → generates data and sends it into the channel
//   Consumer → receives data from the channel and processes it
//
//   Both run in separate goroutines so neither blocks the other.
//   The channel is the handoff point between them.
//
// ============================================================
// KEEPING MAIN ALIVE (learning crutch)
// ============================================================
//
//   When main() returns, ALL running goroutines are killed immediately.
//   time.Sleep() is a quick hack to keep main alive for demos.
//
//   In real programs, use one of:
//   - sync.WaitGroup  → wait for a known number of goroutines to finish
//   - Channel sync    → have goroutine send a "done" signal on a channel
//   - context.Context → cancellation, timeouts, and propagation
//
// ============================================================
// KEY CONCEPTS
// ============================================================
//
//   Goroutine     → lightweight concurrent function (not an OS thread)
//   Channel       → communication mechanism between goroutines
//   Unbuffered    → synchronous; sender and receiver must rendezvous
//   Buffered      → make(chan T, N); sender can send up to N values
//                   without a receiver waiting (introduces a queue)
//   Blocking      → goroutine pauses until the operation can complete;
//                   does NOT freeze the whole program
//   Deadlock      → all goroutines blocked forever; runtime panics
//   Producer      → sends data into a channel
//   Consumer      → receives data from a channel
//
// ============================================================

func CHANNELS() {

	greetings := make(chan string)

	// ----------------------------------------------------------
	// Direct send / receive between two goroutines
	// ----------------------------------------------------------
	// Sender goroutine: sends two values into the channel.
	// Each send blocks until the receiver picks it up.

	go func() {
		greetings <- "Hello Channel"
		greetings <- "World"
	}()

	// Receiver goroutine: receives both values and prints them.
	// Each receive blocks until the sender sends.

	go func() {
		fmt.Println(<-greetings)
		fmt.Println(<-greetings)
	}()

	// ----------------------------------------------------------
	// Producer / Consumer
	// ----------------------------------------------------------
	// range over a string yields runes (Unicode code points),
	// so string(e) converts each rune back to a single-char string.

	// Producer: sends one message per character in "abcde"
	go func() {
		for _, e := range "abcde" {
			greetings <- "Alphabet: " + string(e)
		}
	}()

	// Consumer: receives exactly 5 messages (matches producer output)
	go func() {
		for range 5 {
			fmt.Println(<-greetings)
		}
	}()

	// ----------------------------------------------------------
	// Keep main alive long enough for goroutines to finish.
	// Replace with sync.WaitGroup or channel sync in real code.
	// ----------------------------------------------------------
	time.Sleep(1 * time.Second)

	fmt.Println("End of Program")
}