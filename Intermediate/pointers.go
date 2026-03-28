package main

import "fmt"

func main() {

	var a int = 10

	var ptr *int = &a

	fmt.Println(a)
	fmt.Println(ptr)

	fmt.Println(*ptr)

	if ptr == nil {
		fmt.Println("Pointer is NIL")
	} else {
		fmt.Println("Pointer contains value:", *ptr)
	}

	modifyValue(ptr)

	fmt.Println(a)
}

func modifyValue(ptr *int) {
	*ptr += 1
}


// ============================================================
// Pointer Definition
// ============================================================

// A pointer is a variable that stores the memory address
// of another variable.


// ============================================================
// Pointer Syntax
// ============================================================

// Declaration:
// var ptr *int

// ptr can store the address of an integer variable.


// ============================================================
// Address Operator (&)
// ============================================================

// &variable → returns the memory address of the variable

// Example:
// ptr := &a


// ============================================================
// Dereferencing Operator (*)
// ============================================================

// *ptr → accesses the value stored at the memory address

// Example:
// *ptr gives the value stored at address of a


// ============================================================
// Nil Pointer
// ============================================================

// A pointer that does not point to any memory location
// is called a nil pointer.

// Default value of a pointer = nil.


// ============================================================
// Passing Pointers to Functions
// ============================================================

// When we pass a pointer to a function,
// the function receives the memory address.

// This allows the function to modify the
// original variable directly.


// ============================================================
// Example Flow
// ============================================================

// a = 10

// ptr → stores address of a

// modifyValue(ptr)
// modifies value at that address

// a becomes 11


// ============================================================
// Important Notes
// ============================================================

// Pointers allow indirect modification of variables.

// Go does NOT allow pointer arithmetic
// (unlike C/C++).

// Pointer operations are limited to:
// &  → address-of operator
// *  → dereference operator


// ============================================================
// Unsafe Pointer
// ============================================================

// Go provides unsafe.Pointer for low-level memory access.
// It should be used carefully and rarely.

// Example:
// unsafe.Pointer(&x)


// ============================================================
// Pointer to Pointer
// ============================================================

// It is possible to create pointers to pointers,
// but it is rarely needed in Go.