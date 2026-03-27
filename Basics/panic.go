package main

import "fmt"

func Panic() {

	// ============================================================
	// 🔹 1. Panic Basics
	// ============================================================

	// panic stops normal execution of the current function immediately
	// It is used for unexpected errors where program cannot proceed

	// Syntax:
	// panic(value) → value can be of any type (string, error, etc.)

	var a int
	fmt.Scanln(&a)

	// Process(a)
	NewProcess(a)
}


// ============================================================
// 🔹 2. Simple Panic Example
// ============================================================

func Process(input int) {
	if input < 0 {
		panic("Input must be non-negative number")
	} else {
		fmt.Println("Processing Input")
	}
}


// ============================================================
// 🔹 3. Panic with Defer
// ============================================================

func NewProcess(input int) {

	defer fmt.Println("Function Executed")
	// defer WILL execute even if panic occurs

	if input < 0 {
		fmt.Println("Before Panic")
		panic("Enter Valid Input")
	}

	fmt.Println("Processing Input")
}


// ============================================================
// 🔹 4. Important Notes
// ============================================================

// - panic stops execution immediately (no further normal code runs)
// - Code AFTER panic is NOT executed

// - HOWEVER:
//   defer functions WILL still execute before program crashes

// Execution flow when input < 0:
// 1. "Before Panic" prints
// 2. defer runs → "Function Executed"
// 3. panic message printed → program exits

// ------------------------------------------------------------
// 🔹 Key Rule
// ------------------------------------------------------------

// panic → stops normal flow
// defer → always runs before function exits (even during panic)