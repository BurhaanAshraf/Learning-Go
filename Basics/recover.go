package main

import (
	"fmt"
)

func Recover() {

	// ============================================================
	// 🔹 1. Recover Basics
	// ============================================================

	newProcess()
	fmt.Println("Process Ended")

	// NOTE:
	// - recover is used to handle panic and prevent program crash
	// - It regains control of a panicking function
	// - If panic is recovered → program continues execution
}


// ============================================================
// 🔹 2. Recover with Defer
// ============================================================

func newProcess() {

	defer func() {

		// recover MUST be called inside defer
		// It captures panic value if panic occurred

		r := recover()

		if r != nil {
			fmt.Println("Recovered,", r)
		}

		// NOTE:
		// - If no panic → recover() returns nil
		// - If panic → recover() returns panic value
	}()

	fmt.Println("Start Process")

	panic("Something Went Wrong")
	// panic stops normal execution

	// fmt.Println("End Process") → NOT executed
}


// ============================================================
// 🔹 3. Important Concepts
// ============================================================

// - recover stops panic propagation
// - It ONLY works inside deferred functions
// - It prevents program from crashing
// - After recover → execution continues normally

// Execution Flow:
// 1. "Start Process"
// 2. panic occurs
// 3. deferred function runs
// 4. recover() captures panic
// 5. "Recovered, Something Went Wrong"
// 6. returns to caller
// 7. "Process Ended"


// ------------------------------------------------------------
// 🔹 Key Rule
// ------------------------------------------------------------

// panic   → triggers crash
// defer   → always runs
// recover → stops crash (if used inside defer)