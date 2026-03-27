package main

import "fmt"

func Range() {

	// ============================================================
	// 🔹 1. range over String
	// ============================================================

	message := "Hello, World!"

	for i, v := range message {
		// i → byte index
		// v → rune (Unicode code point)

		fmt.Printf("Index: %d, Rune: %c\n", i, v)

		// NOTE:
		// - Strings are UTF-8 encoded
		// - range iterates over runes (characters)
		// - i is byte index (important for multi-byte chars)
	}

	// Strings are immutable (cannot modify)


	// ============================================================
	// 🔹 2. range over Array
	// ============================================================

	numbers := [5]int{10, 20, 30, 40, 50}

	for i, v := range numbers {
		// i → index
		// v → value
		fmt.Printf("Array Index: %d, Value: %d\n", i, v)
	}

	// NOTE:
	// - Arrays are ordered
	// - Iterates from first to last


	// ============================================================
	// 🔹 3. range over Slice
	// ============================================================

	sliceNumbers := []int{1, 2, 3, 4, 5}

	for i, v := range sliceNumbers {
		fmt.Printf("Slice Index: %d, Value: %d\n", i, v)
	}

	// NOTE:
	// - Same behavior as arrays
	// - Ordered iteration


	// ============================================================
	// 🔹 4. range over Map
	// ============================================================

	studentMarks := map[string]int{
		"Burhaan": 90,
		"Danish":  85,
		"Damir":   88,
	}

	for k, v := range studentMarks {
		// k → key
		// v → value
		fmt.Printf("Key: %s, Value: %d\n", k, v)
	}

	// NOTE:
	// - Maps are unordered
	// - Order may change every iteration


	// ============================================================
	// 🔹 5. Blank Identifier (_)
	// ============================================================

	// Ignore index:
	// for _, v := range sliceNumbers

	// Ignore value:
	// for i := range sliceNumbers

	// NOTE:
	// - "_" discards unused values
}