package main

import "fmt"

func Init() {

	fmt.Println("This is main function")

}


// ============================================================
// 🔹 1. init() Function Basics
// ============================================================

// init() is a special function in Go
// It runs automatically BEFORE main()

// Rules:
// - No parameters
// - No return values
// - Cannot be called manually
// - Executed automatically by Go runtime


// ============================================================
// 🔹 2. Multiple init Functions
// ============================================================

// A file can have multiple init() functions
// They execute sequentially in the order they appear

// Execution order in this file:
// init 01 → init 02 → init 03 → main()


func init() {
	fmt.Println("This is init 01 function")
}

func init() {
	fmt.Println("This is init 02 function")
}

func init() {
	fmt.Println("This is init 03 function")
}


// ============================================================
// 🔹 3. Important Notes
// ============================================================

// - init() runs once per package initialization
// - Used for setup tasks like:
//   • initializing variables
//   • loading configuration
//   • database connections
//   • registering components

// - init() functions run BEFORE main()
// - If multiple files exist in same package,
//   init() execution order depends on file compilation order