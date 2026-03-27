package main

import (
	"fmt"
	"os"
)

func main() {

	// ============================================================
	// 🔹 1. Defer vs os.Exit
	// ============================================================

	defer fmt.Println("Deferred statement")
	// This will NOT run if os.Exit is called

	fmt.Println("Starting the main function")


	// ============================================================
	// 🔹 2. os.Exit
	// ============================================================

	// os.Exit terminates the program IMMEDIATELY
	// No deferred functions are executed

	// Syntax:
	// os.Exit(statusCode)

	// Status Codes:
	// 0 → success
	// non-zero → error / failure

	os.Exit(0)


	// ============================================================
	// 🔹 3. Unreachable Code
	// ============================================================

	// This will NEVER execute
	fmt.Println("End of main function")
}


// ============================================================
// 🔹 4. Important Concepts
// ============================================================

// - os.Exit stops program instantly (no cleanup)
// - defer does NOT run when os.Exit is used
// - Used for critical termination scenarios
// - Status code is returned to OS

// ------------------------------------------------------------
// 🔹 Key Rule
// ------------------------------------------------------------

// defer → runs on normal return or panic
// os.Exit → skips defer and exits immediately