package main

import "fmt"

/*
=============================================================================
CHANNEL CLOSING — REVISION NOTES
=============================================================================

─────────────────────────────────────────────────────────────────────────────
HOW close(ch) WORKS
─────────────────────────────────────────────────────────────────────────────

close(ch) sends a signal that no more values will ever be sent on this channel.
It does NOT delete the channel or the values already in it.

After close:
    → Values already in buffer can still be received normally
    → Once buffer is empty, any receive returns zero value + ok = false instantly
    → Closed channel never blocks on receive — it always returns immediately
    → Sending to a closed channel causes RUNTIME PANIC

─────────────────────────────────────────────────────────────────────────────
GOLDEN RULES OF CLOSING
─────────────────────────────────────────────────────────────────────────────

1. SENDER closes, never receiver.
   The sender knows when it is done sending. The receiver does not.
   If receiver closes and sender tries to send → runtime panic.

2. Close ONCE. Never twice.
   Closing an already closed channel → runtime panic.
   If multiple goroutines might close the same channel, use sync.Once.

3. Never close a nil channel.
   Closing a nil channel → runtime panic.

4. You do not HAVE to close every channel.
   Close is only needed when the receiver needs to know "no more values coming"
   — specifically for (for range) to terminate. If you are receiving manually
   with <-ch and you know how many values to expect, closing is optional.

─────────────────────────────────────────────────────────────────────────────
THREE WAYS TO RECEIVE AND HANDLE CLOSE
─────────────────────────────────────────────────────────────────────────────

WAY 1 — comma-ok (manual, one receive at a time):

    val, ok := <-ch
    if !ok {
        channel is closed and empty
    }

WAY 2 — for range (automatic, drains until closed):

    for val := range ch {
        loops until ch is closed and empty, then stops automatically
    }

WAY 3 — receive from already closed empty channel:

    close(ch)
    val, ok := <-ch   // ok = false, val = zero value, returns instantly

─────────────────────────────────────────────────────────────────────────────
DIRECTIONAL CHANNELS — ENFORCING WHO SENDS AND RECEIVES
─────────────────────────────────────────────────────────────────────────────

When passing channels to functions, restrict direction explicitly:

    chan<- int   →  send-only  (function can only send, cannot receive)
    <-chan int   →  receive-only (function can only receive, cannot send)
    chan int     →  bidirectional (can do both — used when creating channel)

Why this matters:
    → Compiler enforces it — trying to receive on chan<- or send on <-chan
      is a compile error, not a runtime panic.
    → Makes intent clear — reader immediately knows if a function is a
      producer or consumer just from the signature.
    → Only the send-only side can call close(ch). Receiver cannot close
      a receive-only channel — compiler prevents it.

─────────────────────────────────────────────────────────────────────────────
PIPELINE PATTERN — CHAINING CHANNELS THROUGH FUNCTIONS
─────────────────────────────────────────────────────────────────────────────

A pipeline connects functions via channels where output of one is input
of the next. Each stage runs as its own goroutine concurrently.

    main → values(in) → filter(in, out) → main reads out

Each stage:
    → Receives from its input channel (range drains it automatically)
    → Does its work on each value
    → Sends results to its output channel
    → Closes its output channel when input is exhausted

close propagates through the pipeline automatically:
    values closes in → filter's for range stops → filter closes out
    → main's for range stops → program ends cleanly

Without close at each stage, for range in the next stage would block
forever waiting for more values → deadlock.

─────────────────────────────────────────────────────────────────────────────
signal CHANNEL PATTERN — WAITING FOR A GOROUTINE THAT HAS NO RETURN VALUE
─────────────────────────────────────────────────────────────────────────────

When a goroutine does work and then needs to signal completion to main:

    signal := make(chan int)

    go func() {
        do work — for range, processing, etc.
        signal <- 0  // done — value does not matter, only arrival matters
    }()

    <-signal  // main blocks here until goroutine signals completion

This is the done channel pattern applied to a goroutine that uses
for range internally. Main cannot use for range on the same channel
the goroutine is reading — so signal is the coordination mechanism.

─────────────────────────────────────────────────────────────────────────────
BUFFERED vs UNBUFFERED IN THESE PATTERNS
─────────────────────────────────────────────────────────────────────────────

All channels here are unbuffered — sender and receiver must meet.
This means every send blocks until the receiver is ready and vice versa.
That is why goroutines are used — without them, the send on main would
block forever with nobody to receive.

With buffered channels:
    → Sender can send without a receiver being ready (up to buffer capacity)
    → go keyword before values and filter would not be strictly necessary
      for sending — but goroutines are still good practice for concurrency

─────────────────────────────────────────────────────────────────────────────
RUNTIME PANICS TO REMEMBER
─────────────────────────────────────────────────────────────────────────────

    send on closed channel     → panic: send on closed channel
    close of closed channel    → panic: close of closed channel
    close of nil channel       → panic: close of nil channel

All three are runtime panics — not compile errors. The compiler will not
catch them. You must design your channel ownership carefully to avoid them.

─────────────────────────────────────────────────────────────────────────────
SUMMARY — KEY RULES
─────────────────────────────────────────────────────────────────────────────

1.  SENDER closes, never receiver.
2.  Close ONCE — closing twice panics at runtime.
3.  Sending to closed channel panics at runtime.
4.  Receiving from closed empty channel returns zero value + ok=false instantly.
5.  for range stops automatically when channel is closed and drained.
6.  close is only required when receiver needs to know "no more values coming".
7.  chan<- send-only, <-chan receive-only — compiler enforces, use in functions.
8.  Only send-only side can close — receiver cannot close a <-chan channel.
9.  Pipeline: each stage closes its output when its input is exhausted.
10. close propagates through pipeline automatically via for range termination.

=============================================================================
*/

func ClosingChannels() {

	simpleClosing()

	receivingFromClosedChannel()
	
	rangeOverClosedChannel()

	// =========================================================================
	// PIPELINE — values → filter → main
	// =========================================================================
	// Two unbuffered channels connect three stages:
	//   in  → values sends into it, filter reads from it
	//   out → filter sends into it, main reads from it
	//
	// Both functions run as goroutines so they can block on their channel
	// operations concurrently without freezing main.
	//
	// close propagates automatically:
	//   values closes in  → filter's for range stops → filter closes out
	//   → main's for range stops → program ends
	//
	// NOTE: if buffered channels were used, go keyword before values and
	// filter would not be strictly necessary for the sends — buffered channels
	// do not need an immediate receiver. But goroutines are still used here
	// for proper concurrent execution of each pipeline stage.
	// =========================================================================

	in := make(chan int)
	out := make(chan int)

	go values(in)      // producer: generates values, closes in when done
	go filter(in, out) // transformer: reads in, filters, writes out, closes out

	// main is the final consumer — reads filtered values until out is closed
	for val := range out {
		fmt.Println("Filtered Values:", val)
	}
}

// =========================================================================
// EXAMPLE 1 — SIMPLE CHANNEL CLOSING
// =========================================================================
// Goroutine sends 5 values then closes the channel.
// for range on main drains all values and stops when channel is closed.
// Without close(ch), for range would block forever after the 5th value.
// =========================================================================

func simpleClosing() {
	ch := make(chan int)

	go func() {
		for i := range 5 {
			ch <- i
		}
		close(ch) // sender closes — signals no more values coming
		          // for range on main will stop after receiving all 5
	}()

	// stops automatically when ch is closed and drained
	for val := range ch {
		fmt.Println("Received:", val)
	}
}

// =========================================================================
// EXAMPLE 2 — RECEIVING FROM AN ALREADY CLOSED CHANNEL
// =========================================================================
// Closing a channel before anyone receives from it is valid.
// Receive on a closed empty channel returns zero value + ok = false instantly.
// It does NOT block — closed channels always return immediately.
//
// PANIC CASES (do not do these):
//   close(ch) called twice         → panic: close of closed channel
//   ch <- value after close(ch)    → panic: send on closed channel
//   close of nil channel           → panic: close of nil channel
// =========================================================================

func receivingFromClosedChannel() {
	ch := make(chan int)
	close(ch) // close before any send — valid, no panic

	// receive on closed empty channel:
	// val = 0 (zero value for int), ok = false (channel closed and empty)
	// returns instantly — does not block
	val, ok := <-ch

	if !ok {
		fmt.Println("Channel Is Closed") // always reaches here
		return
	}
	fmt.Println("Received:", val) // never reaches here
}

// =========================================================================
// EXAMPLE 3 — for range IN GOROUTINE + signal CHANNEL
// =========================================================================
// Main sends 5 values into ch, then closes ch.
// Goroutine uses for range to drain ch — stops when ch is closed.
// signal channel coordinates: main cannot exit until goroutine finishes.
//
// WHY signal IS NEEDED:
//   Main closes ch and immediately hits <-signal — blocks there.
//   Goroutine finishes its for range (ch closed and drained), sends to signal.
//   Main unblocks and exits cleanly.
//   Without <-signal, main would exit immediately after close(ch),
//   killing the goroutine before it finishes printing all values.
//
// WHY main cannot use for range on ch directly here:
//   Main is the SENDER — it is sending into ch.
//   Goroutine is the RECEIVER — it reads from ch via for range.
//   Sender and receiver must be in different goroutines for unbuffered channels.
// =========================================================================

func rangeOverClosedChannel() {
	ch := make(chan int)
	signal := make(chan int)

	go func() {
		for val := range ch {
			// loops until ch is closed and drained, then stops automatically
			fmt.Println("Values:", val)
		}
		signal <- 0 // goroutine finished all values — signal completion to main
		            // value does not matter, only the arrival matters
	}()

	for i := range 5 {
		ch <- i // each send blocks until goroutine receives — unbuffered rendezvous
	}
	close(ch) // signal goroutine's for range to stop after draining

	<-signal // main blocks here until goroutine confirms it is fully done
	         // without this, main exits and kills goroutine mid-execution
}

// =========================================================================
// PIPELINE STAGE 1 — PRODUCER
// =========================================================================
// Directional channel: chan<- int (send-only)
//   → this function can only send to in, cannot receive from it
//   → compiler enforces this — attempting <-in here is a compile error
//   → only the send-only side can call close — receiver cannot close <-chan
//
// Sends values 0-4 into in, then closes in.
// Closing in signals filter's for range to stop after draining.
// =========================================================================

func values(in chan<- int) {
	for i := range 5 {
		in <- i // send each value — blocks until filter receives
	}
	close(in) // sender closes — filter's for range will terminate after this
}

// =========================================================================
// PIPELINE STAGE 2 — TRANSFORMER / FILTER
// =========================================================================
// Two directional channels:
//   in  → <-chan int (receive-only) — can only read from in
//   out → chan<- int (send-only)    — can only send to out
//
// for range on in drains all values until values() closes in.
// Only even numbers pass through to out.
// After in is exhausted, close(out) propagates the termination signal
// forward — main's for range on out will stop after this.
//
// This is how close propagates through a pipeline:
//   values closes in → filter's range stops → filter closes out
//   → main's range stops → program ends cleanly
// =========================================================================

func filter(in <-chan int, out chan<- int) {
	for i := range in {
		// for range receives from in until in is closed and drained
		if i%2 == 0 {
			out <- i // only even numbers go into out
		}
	}
	close(out) // in exhausted — close out so main's for range terminates
}