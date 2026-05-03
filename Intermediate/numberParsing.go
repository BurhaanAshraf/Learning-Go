package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func NumberParsing() {

	// ============================================================
	// 🔹 String → Integer
	// ============================================================

	// Atoi()
	// ASCII to Integer

	numStr := "12345"

	intVal, err := strconv.Atoi(numStr)

	if err != nil {
		panic("Something Went Wrong While Parsing")
	}

	fmt.Println("Parsed Integer:", intVal)

	// %T prints type
	fmt.Printf("Type: %T\n", intVal)


	// ============================================================
	// 🔹 ParseInt()
	// ============================================================

	// ParseInt(string, base, bitSize)

	// base:
	// 2  → binary
	// 10 → decimal
	// 16 → hexadecimal

	// bitSize:
	// 8,16,32,64

	newStr := "123456789"

	parsedInt, err := strconv.ParseInt(
		newStr,
		10,
		32,
	)

	if err != nil {
		panic("Something Went Wrong")
	}

	fmt.Println("Parsed Int:", parsedInt)

	// ParseInt always returns int64
	fmt.Printf("Type: %T\n", parsedInt)


	// ============================================================
	// 🔹 String → Float
	// ============================================================

	// ParseFloat(string, bitSize)

	floatStr := "3.1415"

	floatVal, err := strconv.ParseFloat(
		floatStr,
		64,
	)

	if err != nil {
		panic("Something Went Wrong With Parsing")
	}

	// %.2f → 2 decimal places
	fmt.Printf("Float Value: %.2f\n", floatVal)

	// reflect.TypeOf() checks type
	fmt.Println(
		"Float Type:",
		reflect.TypeOf(floatVal),
	)


	// ============================================================
	// 🔹 Binary → Decimal
	// ============================================================

	binaryDigits := "1010"

	// base 2 means binary parsing
	decimalVal, err := strconv.ParseInt(
		binaryDigits,
		2,
		64,
	)

	if err != nil {
		panic("Something Went Wrong With Parsing")
	}

	fmt.Println(
		"Binary to Decimal:",
		decimalVal,
	)


	// ============================================================
	// 🔹 Parsing Error Example
	// ============================================================

	invalidNum := "1234abc"

	invalidParse, err := strconv.Atoi(invalidNum)

	if err != nil {

		// strconv returns parsing error
		fmt.Println(
			"Parsing Error:",
			err,
		)

		return
	}

	fmt.Println(invalidParse)
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// strconv.Atoi()
// → string to int

// strconv.ParseInt()
// → string to int64

// ParseInt always returns int64
// regardless of bitSize passed
// bitSize only limits parsing range

// strconv.ParseFloat()
// → string to float

// base:
// 2 → binary
// 10 → decimal
// 16 → hexadecimal

// ParseInt always returns int64

// reflect.TypeOf()
// → check variable type

// %.2f
// → 2 decimal places