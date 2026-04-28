package main

import (
	"errors"
	"fmt"
)

// ============================================================
// 🔹 Custom Error
// ============================================================

// Custom errors help us store:
// code + message + original error

type customError struct {
	code    int
	message string
	er      error
}


// ============================================================
// 🔹 Error() Method
// ============================================================

// Implement Error() string
// to satisfy built-in error interface

func (err *customError) Error() string {
	return fmt.Sprintf(
		"Error Code %d: %s, %v",
		err.code,
		err.message,
		err.er,
	)
}


// ============================================================
// 🔹 High-Level Function
// ============================================================

// High-level = business logic
// calls lower-level function
// adds more context

func doSomething() error {

	err := doSomethingElse()

	if err != nil {
		return &customError{
			code:    500,
			message: "Something went wrong!",
			er:      err,
		}
	}

	return nil
}


// ============================================================
// 🔹 Low-Level Function
// ============================================================

// Low-level = actual internal work
// returns original error

func doSomethingElse() error {
	return errors.New("Internal Error")
}


// ============================================================
// 🔹 Main Function
// ============================================================

func CustomError() {

	err := doSomething()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Operation Completed!")
}


// ============================================================
// 🔹 Quick Revision
// ============================================================

// Error() → required for custom errors

// Low-level → original error

// High-level → adds context

// Better debugging + clean code