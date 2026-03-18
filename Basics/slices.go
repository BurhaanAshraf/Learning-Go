package main

import (
	"fmt"
	"slices"
)

func slice() {

	// ============================================================
	// 🔹 1. Declaring & Initializing Slices
	// ============================================================

	// Basic slice declaration
	var numbers []int = []int{1, 2, 3, 4, 5, 6}
	_ = numbers

	// Short-hand declaration (idiomatic Go)
	floatingNumbers := []int{1, 2, 3, 4}
	_ = floatingNumbers

	// Slice of strings
	var names []string = []string{"Burhaan", "Danish", "Damir"}
	_ = names

	// NOTE:
	// - Slices are dynamic (no fixed size like arrays)
	// - Most common syntax: nums := []int{1,2,3}


	// ============================================================
	// 🔹 2. Creating Slice using make()
	// ============================================================

	slice := make([]int, 5)
	_ = slice

	// NOTE:
	// - make([]type, length)
	// - Initializes with zero values
	// - length = capacity (by default)


	// ============================================================
	// 🔹 3. Slicing (Creating Sub-slices)
	// ============================================================

	a := []int{1, 2, 3, 4, 5, 6}

	slice1 := a[1:4] // from index 1 to index before 4
	fmt.Println(slice1) // Output: [2 3 4]

	// NOTE:
	// - Syntax: slice[start:end]
	// - start included, end excluded


	// ============================================================
	// 🔹 4. Appending to a Slice
	// ============================================================

	slice1 = append(slice1, 6, 7)
	fmt.Println("Slices:", slice1)

	// NOTE:
	// - append returns a new slice
	// - underlying array may change if capacity exceeded


	// ============================================================
	// 🔹 5. Copying a Slice
	// ============================================================

	sliceCopy := make([]int, len(slice1))

	// copy(destination, source)
	copy(sliceCopy, slice1)

	fmt.Println("Slices Copy:", sliceCopy)

	// NOTE:
	// - copy creates a deep copy (different underlying array)


	// ============================================================
	// 🔹 6. Nil Slice
	// ============================================================

	var nilSlice []int // no underlying array
	_ = nilSlice

	// NOTE:
	// - nil slice == no memory allocated
	// - len = 0, cap = 0
	// - useful in many Go patterns


	// ============================================================
	// 🔹 7. Iterating over a Slice
	// ============================================================

	for i, v := range slice1 {
		fmt.Printf("Index %d Value %d\n", i, v)
	}

	// NOTE:
	// - range returns index + value


	// ============================================================
	// 🔹 8. Modifying Slice Values
	// ============================================================

	slice1[0] = 100

	copy(sliceCopy, slice1)

	if slices.Equal(slice1, sliceCopy) {
		fmt.Println("Equal")
	} else {
		fmt.Println("Not Equal")
	}

	// NOTE:
	// - slices are reference-like (share underlying array)


	// ============================================================
	// 🔹 9. 2D Slices
	// ============================================================

	twoDSlice := make([][]int, 3)

	for i := 0; i < len(twoDSlice); i++ {

		twoDSlice[i] = make([]int, 3)

		for j := 0; j < 3; j++ {
			twoDSlice[i][j] = i + j*10
		}
	}

	fmt.Println(twoDSlice)

	// NOTE:
	// - slice of slices (like dynamic 2D array)
	// - each inner slice is created separately


	// ============================================================
	// 🔹 10. Length vs Capacity
	// ============================================================

	checking := []int{1, 2, 3, 4, 5, 6, 7, 8}

	newChecking := checking[2:5]

	// Length:
	// elements from index 2 to before 5 → [3 4 5] → length = 3

	// Capacity:
	// from index 2 till end → [3 4 5 6 7 8] → capacity = 6

	fmt.Printf("Length of checking is %d , Capacity of checking is %d\n",
		len(checking), cap(checking))

	fmt.Printf("Length of newChecking is %d , Capacity of newChecking is %d\n",
		len(newChecking), cap(newChecking))
}
// ============================================================
// 🔹 IMPORTANT: Slices share underlying array (CORE CONCEPT)
// ============================================================

// A slice does NOT store data itself.
// It is just a "view" (or window) over an underlying array.

// Internally, a slice contains:
// - pointer → points to array
// - length  → number of elements
// - capacity → max usable elements ahead

// ------------------------------------------------------------
// 🔸 Example of sharing memory
// ------------------------------------------------------------

// a := []int{1, 2, 3, 4}
// b := a[1:3]   // b = [2 3]

// Both a and b point to SAME underlying array:
// [1 2 3 4]
//     ↑ ↑
//     b b

// ------------------------------------------------------------
// 🔸 Modifying one affects the other
// ------------------------------------------------------------

// b[0] = 100

// Now:
// a = [1 100 3 4]
// b = [100 3]

// Reason:
// → both slices share the SAME memory

// ------------------------------------------------------------
// 🔸 Why this can be dangerous
// ------------------------------------------------------------

// original := []int{1, 2, 3, 4, 5}
// sub := original[1:4]

// sub[0] = 999

// Now:
// original = [1 999 3 4 5]  ❗ (unexpected change)

// ------------------------------------------------------------
// 🔸 Solution: Use copy() for independent data
// ------------------------------------------------------------

// safe := make([]int, len(sub))
// copy(safe, sub)

// safe[0] = 111

// Now:
// original = [1 999 3 4 5]
// safe     = [111 3 4]  ✅ independent

// ------------------------------------------------------------
// 🔸 Special Case: append()
// ------------------------------------------------------------

// append may OR may not create new memory:

// If capacity is NOT exceeded:
// → still shares underlying array

// If capacity IS exceeded:
// → Go creates NEW array (no sharing)

// ------------------------------------------------------------
// 🔸 Golden Rule
// ------------------------------------------------------------

// Slicing  → shared memory
// copy()   → new memory (safe)
// append() → depends on capacity

// ------------------------------------------------------------
// 🔸 One-line intuition
// ------------------------------------------------------------

// "Slice is a lens over an array, not the actual data"