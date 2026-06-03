package main

import (
	"fmt"
	"time"
)

func unbufferedChannels() {

    // buffered channels have internal storage (a queue of fixed capacity)

    // main difference between buffered and unbuffered channels:
    // buffered channels allow asynchronous communication — the sender can
    // continue working without blocking, as long as the buffer is not full.
    // unbuffered channels require sender and receiver to both be ready at
    // the exact same moment (synchronous handoff / rendezvous).

    // buffered channels let us control the rate of data flow between
    // producers and consumers by tuning the buffer size

    greetings := make(chan int)

    go func() {
        rcvr := <-greetings
        fmt.Println(rcvr)
        time.Sleep(2 * time.Second)
        // goroutine 1 is the receiver. it immediately blocks here waiting for
        // a value. once goroutine 2 sends the value, both goroutines unblock
        // instantly — it is a direct handoff between goroutine 2 and goroutine 1.
        // main is just sleeping at this point; it plays no role in the handoff.
        fmt.Println("2 seconds passed")
    }()

    go func() {
        greetings <- 2
        // goroutine 2 is the sender. it sends 2 directly to goroutine 1 (the receiver).
        // this is the only send in the program — the channel now has no more values.
    }()

    go func() {
        time.Sleep(3 * time.Second)
        fmt.Println("3 seconds passed")
    }()

    time.Sleep(4 * time.Second)

    rcvr := <-greetings
    fmt.Println(rcvr)

    fmt.Println("End of Program") // ← never reached — intentional deadlock below explains why

    // DEADLOCK DEMO:
    // by the time main wakes from its 4s sleep, goroutines 1, 2, and 3 have all
    // already finished. the channel is empty and nobody is left to send a value.
    // main blocks forever on "rcvr := <-greetings" waiting for a value that
    // will never arrive → Go runtime detects all goroutines are asleep → deadlock!

    // KEY RULE for unbuffered channels:
    // every send must have a matching receive, and every receive must have a
    // matching send — happening concurrently (in different goroutines).
    // if one side is missing, the waiting side blocks forever → deadlock.

    // goroutines are used with unbuffered channels precisely for this reason:
    // they let sender and receiver run concurrently so neither blocks the other
    // permanently. without goroutines, you would try to send and receive on the
    // same goroutine sequentially — but the first operation would block forever,
    // and you would never reach the second one.

    // NOTE: receivers CAN be inside goroutines (goroutine 1 above proves this).
    // what matters is not WHERE the receiver is, but that a matching sender
    // exists and runs concurrently. here, main itself is the receiver and
    // no matching sender exists anymore — that is why it deadlocks.

    // if there is no goroutine alive to send a value to a waiting receiver
    // (or to receive from a waiting sender), Go detects it and panics with
    // "fatal error: all goroutines are asleep — deadlock!"
}