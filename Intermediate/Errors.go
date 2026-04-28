package main

import (
	"errors"
	"fmt"
)

// ============================================================
// 🔹 Error Definition
// ============================================================

// In Go, errors are handled by returning them as values.

// Functions usually return:
// result, error

// If error == nil → success
// If error != nil → something went wrong

// Go does not use try-catch like Java

// ============================================================
// 🔹 Built-in error Interface
// ============================================================

// error is a built-in interface in Go

// type error interface {
//     Error() string
// }

// Any type that implements:
//
// Error() string
//
// automatically satisfies the error interface

// ============================================================
// 🔹 Basic Error Example
// ============================================================

func sqrt(x float64) (float64, error) {

	if x < 0 {
		return 0, errors.New(
			"math error: square root of negative number",
		)
	}

	return 1, nil
}


// ============================================================
// 🔹 Another Error Example
// ============================================================

func process(x []byte) error {

	if len(x) == 0 {
		return errors.New("Error: Empty Data")
	}

	return nil
}


// ============================================================
// 🔹 Why We Sometimes Don't Write Error()
// ============================================================

// errors.New()
// and fmt.Errorf()

// already return values whose types
// already implement Error() string

// That is why we do not manually write
// Error() when using them


// ============================================================
// 🔹 Custom Error Type
// ============================================================

// When we create our own custom error struct,
// then we MUST implement Error() string

type myError struct {
	message string
}


// This method makes myError satisfy
// the built-in error interface

func (err *myError) Error() string {
	return fmt.Sprintf("Error: %s", err.message)
}

func processError() error {
	return &myError{
		"Custom Error Message",
	}
}


// ============================================================
// 🔹 Wrapped Errors
// ============================================================

// %w is used to wrap errors
// so original error information is preserved

func readConfig() error {
	return errors.New("Config Error")
}

func readData() error {

	err := readConfig()

	if err != nil {
		return fmt.Errorf("Read Data: %w", err)
	}

	return nil
}


// ============================================================
// 🔹 Main Function
// ============================================================

func Error() {

	res, err := sqrt(16)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)


	// Example with process()

	data := []byte{1, 2, 3}

	newErr := process(data)

	if newErr != nil {
		fmt.Println(newErr)
		return
	}

	fmt.Println("Data processed successfully")


	// Custom Error Example

	customErr := processError()
	fmt.Println(customErr)


	// Wrapped Error Example

	configErr := readData()

	if configErr != nil {
		fmt.Println(configErr)
		return
	}

	fmt.Println("Data Read Successfully")
}


// ============================================================
// 🔹 Quick Revision Points
// ============================================================

// errors.New()
// → creates simple error

// fmt.Errorf()
// → formatted error creation

// %w
// → wraps original error

// Error() string
// → required for custom errors

// err != nil
// → check for failure

// Go errors are values, not exceptions