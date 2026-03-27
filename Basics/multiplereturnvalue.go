package main

import (
	"errors"
	"fmt"
)

func Multiple() {

	// ============================================================
	// 🔹 1. Multiple Return Values
	// ============================================================

	// Go functions can return multiple values
	// Syntax:
	// func name(params) (type1, type2) { return val1, val2 }

	a := 10
	b := 4

	quotient, remainder := divide(a, b)
	fmt.Printf("Quotient: %d , Remainder %d\n", quotient, remainder)


	// ============================================================
	// 🔹 2. Error Handling Pattern (VERY IMPORTANT)
	// ============================================================

	res, err := compare(10, 20)

	if err != nil {
		// If error exists → handle it
		fmt.Println(err)
	} else {
		// Otherwise use result
		fmt.Println(res)
	}

	// NOTE:
	// - Go uses explicit error handling (no exceptions)
	// - Convention: return (value, error)
	// - err == nil → success
	// - err != nil → something went wrong


	// ============================================================
	// 🔹 3. Named Return Values
	// ============================================================

	newQ, newR := Newdivide(10, 6)
	fmt.Printf("NewQ: %v , NewR: %v\n", newQ, newR)

	// NOTE:
	// - Return variables can be named in function signature
	// - "return" alone returns those variables
}


// ============================================================
// 🔹 4. Basic Multiple Return Function
// ============================================================

func divide(a int, b int) (int, int) {

	quotient := a / b
	remainder := a % b

	return quotient, remainder
}


// ============================================================
// 🔹 5. Returning Error with Value
// ============================================================

func compare(a int, b int) (string, error) {

	if a > b {
		return "A is greater than B", nil
		// nil → no error
	} else if b > a {
		return "B is greater than A", nil
	} else {
		return "", errors.New("Unable to compare which is greater")
		// return zero value + error
	}
}

// NOTE:
// - If error occurs → return zero value + error
// - If no error → return actual value + nil


// ============================================================
// 🔹 6. Named Return Values Function
// ============================================================

func Newdivide(a int, b int) (quotient int, remainder int) {

	quotient = a / b
	remainder = a % b

	return
	// Automatically returns quotient and remainder
}