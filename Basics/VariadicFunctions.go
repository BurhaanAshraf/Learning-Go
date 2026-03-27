package main

import "fmt"

func variadic() {

	// ============================================================
	// 🔹 1. Variadic Functions (Basics)
	// ============================================================

	// Variadic functions allow passing variable number of arguments
	// Syntax: ...type  → accepts 0 or more values

	// func name(param1 type1, param2 ...type2) returnType { }

	statement, res := sum("Sum:", 1, 2, 3, 4, 5, 6)
	fmt.Println(statement, res)


	// ============================================================
	// 🔹 2. Passing Slice to Variadic Function
	// ============================================================

	numbers := []int{1, 2, 3, 4, 5}

	fmt.Println(multiply(numbers...))
	// "..." unpacks slice into individual arguments


	// ============================================================
	// 🔹 3. Key Rules
	// ============================================================

	// - Variadic parameter must be LAST parameter
	// - Inside function → variadic param behaves like a slice
	// - Can accept zero or more values
}


// ============================================================
// 🔹 4. Variadic Function Example
// ============================================================

func sum(str string, nums ...int) (string, int) {

	total := 0

	for _, val := range nums {
		// "_" ignores index (blank identifier)
		total += val
	}

	return str, total
}

// NOTE:
// - str → normal parameter
// - nums → variadic parameter (slice of ints)


// ============================================================
// 🔹 5. Another Variadic Function
// ============================================================

func multiply(numbers ...int) int {

	total := 1

	for _, v := range numbers {
		total *= v
	}

	return total
}


// ============================================================
// 🔹 6. Important Concepts
// ============================================================

// - Variadic functions handle dynamic inputs
// - nums ...int → internally treated as []int
// - Use slice... to pass slice as arguments
// - Useful when number of inputs is unknown