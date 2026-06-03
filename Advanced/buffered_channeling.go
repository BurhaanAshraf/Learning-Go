package main

import (
	"fmt"
	"time"
)

/*
=============================================================================
BUFFERED CHANNELS — COMPLETE REVISION NOTES
=============================================================================

─────────────────────────────────────────────────────────────────────────────
MENTAL MODEL — BUFFERED CHANNEL IS A FIFO QUEUE
─────────────────────────────────────────────────────────────────────────────

An unbuffered channel requires sender and receiver to meet at the same time.
A buffered channel removes that requirement by adding a queue in between.
The sender drops the value into the queue and moves on immediately.
The receiver picks it up from the queue whenever it is ready.
They do not need to meet. The queue is the middleman.

    Sender  →  [ slot | slot | slot ]  →  Receiver
                      Buffer (FIFO Queue)

FIFO = First In, First Out.
The first value sent is the first value received. Order is always preserved.

Example:

    ch := make(chan int, 2)     // buffer with 2 slots

    State after creation:
    [ _ | _ ]

    ch <- 10
    [ 10 | _ ]

    ch <- 20
    [ 10 | 20 ]

    <-ch   → returns 10 (first in, first out)
    [ 20 | _ ]

─────────────────────────────────────────────────────────────────────────────
CAPACITY vs LENGTH
─────────────────────────────────────────────────────────────────────────────

    cap(ch)  →  maximum number of values the channel can hold at once
    len(ch)  →  number of values currently sitting in the buffer

    ch := make(chan int, 5)
    ch <- 1
    ch <- 2

    len(ch) = 2   // 2 values currently in buffer
    cap(ch) = 5   // can hold up to 5

    These are useful for debugging and understanding buffer state at any point.

─────────────────────────────────────────────────────────────────────────────
BLOCKING RULES — THE TWO CASES WHERE GOROUTINES BLOCK
─────────────────────────────────────────────────────────────────────────────

SEND  →  ch <- value

    Is there a free slot in the buffer?
    YES → value is stored, goroutine continues immediately
    NO  → goroutine BLOCKS and waits until a receiver frees up a slot

RECEIVE  →  value := <-ch

    Is there a value in the buffer?
    YES → value is returned, goroutine continues immediately
    NO  → goroutine BLOCKS and waits until a sender puts a value in

Summary:
    Buffered channel blocks on SEND  only when buffer is FULL.
    Buffered channel blocks on RECEIVE only when buffer is EMPTY.

─────────────────────────────────────────────────────────────────────────────
CRITICAL DISTINCTION — CHANNELS DO NOT BLOCK. GOROUTINES DO.
─────────────────────────────────────────────────────────────────────────────

The channel itself is just a queue. It has no concept of being blocked.
It is always the goroutine that gets blocked while waiting on the channel.

    WRONG:   "The channel is blocked"
    CORRECT: "The main goroutine is blocked on this channel"
    CORRECT: "The worker goroutine is blocked waiting for a value"

This distinction matters when reasoning about deadlocks and concurrency bugs.

─────────────────────────────────────────────────────────────────────────────
GOROUTINE STATES — BLOCKED vs RUNNABLE vs RUNNING
─────────────────────────────────────────────────────────────────────────────

Every goroutine is always in one of these three states:

    BLOCKED   →  cannot make progress, waiting for something (channel, lock, etc.)
    RUNNABLE  →  ready to run, but waiting for the scheduler to pick it up
    RUNNING   →  currently executing on the CPU

Transition:
    BLOCKED → RUNNABLE → RUNNING

A goroutine does NOT go directly from BLOCKED to RUNNING.
It first becomes RUNNABLE, then the Go scheduler decides when to run it.

─────────────────────────────────────────────────────────────────────────────
COMMON MISTAKE — ASSUMING IMMEDIATE EXECUTION AFTER UNBLOCKING
─────────────────────────────────────────────────────────────────────────────

When a goroutine gets unblocked (sender gets space, receiver gets a value),
it does NOT run immediately. It becomes RUNNABLE and the scheduler decides.

    WRONG mental model:
        send value → receiver immediately runs

    CORRECT mental model:
        send value → receiver becomes RUNNABLE → scheduler decides → receiver RUNS

    WRONG mental model:
        receive value → sender immediately runs

    CORRECT mental model:
        receive value → sender becomes RUNNABLE → scheduler decides → sender RUNS

The current goroutine may keep running for a while before the scheduler
switches to the newly unblocked goroutine. Never assume order of execution
based on when something gets unblocked.

─────────────────────────────────────────────────────────────────────────────
BUFFER FULL SCENARIO — SENDER GETS BLOCKED
─────────────────────────────────────────────────────────────────────────────

    ch := make(chan int, 2)
    ch <- 10
    ch <- 20
     Buffer: [ 10 | 20 ] — full

    ch <- 30
    No free slot → sender goroutine BLOCKS here

     Meanwhile, another goroutine does:
    <-ch
     Buffer: [ 20 ] — one slot freed

    Sender goroutine becomes RUNNABLE
    Scheduler eventually runs it
     Buffer: [ 20 | 30 ]

─────────────────────────────────────────────────────────────────────────────
BUFFER EMPTY SCENARIO — RECEIVER GETS BLOCKED
─────────────────────────────────────────────────────────────────────────────

    ch := make(chan int, 2)
    Buffer: [] — empty

    <-ch
     No value available → receiver goroutine BLOCKS here

    Meanwhile, another goroutine does:
    ch <- 10
    Buffer: [ 10 ]

    Receiver goroutine becomes RUNNABLE
    Scheduler eventually runs it
    Receiver gets 10, buffer: []

─────────────────────────────────────────────────────────────────────────────
DEADLOCK WITH BUFFERED CHANNELS
─────────────────────────────────────────────────────────────────────────────

Deadlock happens when ALL goroutines are blocked and none can make progress.

    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    ch <- 3   // buffer full, main goroutine blocks
               no other goroutine exists to receive
               nobody can free up space
               → DEADLOCK

Deadlock is not about timing. It is about a permanent inability to proceed.
If there is always a goroutine alive that will eventually send or receive,
there is no deadlock — just waiting.

─────────────────────────────────────────────────────────────────────────────
BUFFERED vs UNBUFFERED — SIDE BY SIDE
─────────────────────────────────────────────────────────────────────────────

    Unbuffered:  make(chan int)
    Buffered:    make(chan int, N)

    Unbuffered:
        → No queue. Zero capacity.
        → Sender and receiver MUST meet at the same moment (rendezvous).
        → Think: handshake 🤝 — both hands must be there at the same time.
        → Goroutines are required to make this work without deadlock.

    Buffered:
        → Has a queue of capacity N.
        → Sender does not need a receiver to be waiting — just needs a free slot.
        → Receiver does not need a sender to be waiting — just needs a value in buffer.
        → Decouples sender and receiver in time.
        → Think: mailbox 📬 — sender drops the letter, receiver picks it up later.

─────────────────────────────────────────────────────────────────────────────
REVISION CHECKLIST — HOW TO REASON ABOUT ANY BUFFERED CHANNEL PROGRAM
─────────────────────────────────────────────────────────────────────────────

When you see  ch := make(chan int, N)  ask yourself at each point in the code:

    1. What values are currently in the buffer?
    2. Is the buffer full? (len == cap)
    3. Is the buffer empty? (len == 0)
    4. Which goroutine is trying to send or receive right now?
    5. Will it block or continue?
    6. If it blocks, which goroutine can unblock it?
    7. When that goroutine unblocks the first, what state does it transition to?
    8. What does the scheduler choose to run next?

If you can answer all of these for every line, you can predict any
buffered channel program with confidence.

=============================================================================
*/

// =============================================================================
// EXAMPLE 1 — BLOCKING ON SEND WHEN BUFFER IS FULL
// =============================================================================
//
// Demonstrates: sender goroutine blocks when buffer is at capacity,
// and gets unblocked only after a receiver frees up a slot.

func example1() {

	ch := make(chan int, 2) // buffered channel with capacity 2
	                        // buffer state: [ _ | _ ]

	ch <- 1 // buffer: [ 1 | _ ] — slot used, no blocking, sender continues
	ch <- 2 // buffer: [ 1 | 2 ] — buffer now FULL

	fmt.Println("Receiving From Buffer")

	// This goroutine will sleep 2 seconds then receive from the channel.
	// The sleep simulates a slow consumer. During those 2 seconds, main
	// will try to send ch <- 3 and block because the buffer is full.
	// Once this goroutine wakes and receives, it frees a slot, which
	// makes the main goroutine RUNNABLE again.
	go func() {
		fmt.Println("Goroutine started — will receive after 2 seconds")
		time.Sleep(2 * time.Second)
		// receives 1 (FIFO — first value sent is first value out)
		// buffer after receive: [ 2 | _ ] — one slot freed
		fmt.Println("Received:", <-ch)
	}()

	fmt.Println("Blocking Starts — buffer is full, trying to send ch <- 3")

	// buffer is full [ 1 | 2 ], no free slot.
	// main goroutine BLOCKS here until the goroutine above receives and frees a slot.
	// once slot is freed: main becomes RUNNABLE → scheduler runs it → sends 3.
	// buffer after send: [ 2 | 3 ]
	ch <- 3

	fmt.Println("Blocking Ends — slot was freed, ch <- 3 completed")

	// give the goroutine time to finish printing before main exits.
	// without this sleep, main may exit before the goroutine prints its output.
	time.Sleep(3 * time.Second)

	// KEY TAKEAWAY:
	// buffered channels block on SEND only when the buffer is FULL.
	// buffered channels block on RECEIVE only when the buffer is EMPTY.
	// between those two extremes, sender and receiver are fully independent.
}

// =============================================================================
// EXAMPLE 2 — BLOCKING ON RECEIVE WHEN BUFFER IS EMPTY
// =============================================================================
//
// Demonstrates: receiver goroutine blocks when buffer has no values,
// and gets unblocked only after a sender puts a value in.

func example2() {

	ch := make(chan int, 2) // buffered channel, capacity 2
	                        // buffer state: [ _ | _ ] — empty

	// This goroutine sleeps 2 seconds then sends two values.
	// During those 2 seconds, main will try to receive and block
	// because the buffer is empty — no values to read yet.
	// Once the goroutine sends, main becomes RUNNABLE and reads the values.
	go func() {
		time.Sleep(2 * time.Second)
		ch <- 1 // buffer: [ 1 | _ ]
		ch <- 2 // buffer: [ 1 | 2 ]
	}()

	// buffer is empty right now — main goroutine BLOCKS here.
	// it waits until the goroutine above sends a value into the channel.
	// once a value arrives: main becomes RUNNABLE → scheduler runs it → reads value.
	fmt.Println("Value:", <-ch) // receives 1 (FIFO)
	fmt.Println("Value:", <-ch) // receives 2

	fmt.Println("End of Program")

	// KEY TAKEAWAY:
	// the receiver did not need to wait for the sender to be ready at the same time.
	// it just waited for a value to appear in the buffer.
	// this is what makes buffered channels different from unbuffered —
	// the buffer decouples sender and receiver in time.

}