package main

import "fmt"

func type_conversions() {

	// ------------------------------------------------
	// TYPE CONVERSION
	// ------------------------------------------------
	//
	// Go does NOT perform implicit type conversion.
	//
	// Wrong:
	// var b int32 = a
	//
	// Correct:
	// b := int32(a)
	//
	// Syntax:
	// Type(value)
	//
	var a int = 32

	// int -> int32
	b := int32(a)

	// int32 -> float64
	c := float64(b)

	// Not all conversions are allowed.
	//
	// Example:
	// bool cannot be created from float64.
	//
	// d := bool(c) ❌
	//
	// Go requires explicit and valid conversions.

	e := 3.14

	// float64 -> int
	//
	// Decimal part is discarded,
	// not rounded.
	//
	// 3.14 -> 3
	f := int(e)

	fmt.Println(c, f)

	// ------------------------------------------------
	// STRING TO BYTE SLICE
	// ------------------------------------------------
	//
	// Strings are immutable.
	//
	// Converting to []byte allows working
	// with individual bytes.
	//
	// Commonly used for:
	// - File operations
	// - Network operations
	// - JSON/XML processing
	//
	// Syntax:
	// []byte(string)
	//
	g := "Hello"

	var H []byte

	H = []byte(g)

	fmt.Println(H)

	// ------------------------------------------------
	// BYTE SLICE TO STRING
	// ------------------------------------------------
	//
	// byte is an alias for uint8.
	//
	// Range:
	// 0 - 255
	//
	// Therefore byte slices cannot contain
	// values greater than 255.
	//
	// Common Interview Question:
	// What is byte in Go?
	//
	// Answer:
	// byte is an alias for uint8.
	//
	i := []byte{110, 123, 126}

	// Converts bytes into corresponding
	// ASCII/UTF-8 characters.
	fmt.Println(string(i))

	/*
		Quick Revision

		Type Conversion:
			Type(value)

		Go does not support
		implicit conversions.

		float64 -> int
			Removes decimal part

		byte
			Alias for uint8

		byte range
			0 - 255

		string -> []byte
			[]byte(str)

		[]byte -> string
			string(bytes)
	*/
}