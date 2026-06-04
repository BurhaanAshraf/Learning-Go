package main

import "fmt"

func channelDirs() {
	// NOTES:
	// 1. By default, make(chan int) creates a BIDIRECTIONAL channel.
	//    It can both send data (ch <-) and receive data (<-ch).
	// 2. We keep it bidirectional here in main so it can act as the 
	//    bridge between our producer and consumer.
	ch := make(chan int)

	// 3. Implied Type Conversion:
	//    Go automatically "downgrades" the bidirectional channel 'ch'
	//    to a send-only channel when passed into producer(), and to a
	//    receive-only channel when passed into consumer().
	producer(ch)
	consumer(ch)
}

// Send-Only Channel (chan <- int)
// NOTES:
// 1. Syntax: The arrow points INTO the channel 'chan <-'.
// 2. Purpose: Enforces a strict boundary. This function can ONLY send data.
//    Attempting to read (<-ch) here will cause a compile-time error.
// 3. Compatibility: Accepts bidirectional or send-only channels. 
//    Will NOT accept a receive-only channel.
func producer(ch chan <- int) {
	// Run producer in a separate goroutine so it can send values while the consumer receives them.
	go func() {
		for i := range 5 {
			ch <- i
		}
		// 4. The Ownership Principle:
		//    Always close the channel from the sender's side. 
		//    Since this is a send-only channel, it is legally allowed to close it.
		close(ch)
	}()
}

// Receive-Only Channel (<- chan int)
// NOTES:
// 1. Syntax: The arrow points AWAY from the channel '<- chan'.
// 2. Purpose: Ensures this function can ONLY read data. It cannot inject 
//    accidental data or close the channel (closing a receive-only channel is a compiler error).
// 3. Compatibility: Accepts bidirectional or receive-only channels. 
//    Will NOT accept a send-only channel.
func consumer (ch <- chan int) {
	// 4. The range loop automatically keeps reading from the channel until 
	//    the channel is closed by the producer. Once closed, the loop breaks cleanly.
	for val := range ch {
		fmt.Println("Received: ", val)
	}
}