package main

import (
	"fmt"
	"time"
)

func Epoch() {

	// ============================================================
	// 🔹 Unix / Epoch Time
	// ============================================================

	// Epoch Time (Unix Time)
	// starts from:
	//
	// 00:00:00 UTC
	// January 1, 1970
	//
	// Unix time stores:
	// total seconds passed since this moment

	// Very useful for:
	// - timestamps
	// - databases
	// - APIs
	// - distributed systems
	// - comparing times globally


	// ============================================================
	// 🔹 Current Time
	// ============================================================

	now := time.Now()

	fmt.Println("Current Time:", now)


	// ============================================================
	// 🔹 Convert Current Time → Unix Time
	// ============================================================

	// Unix() converts time → epoch seconds

	unixTimestamp := now.Unix()

	fmt.Println(
		"Current Unix Timestamp:",
		unixTimestamp,
	)


	// ============================================================
	// 🔹 Convert Unix Time → Human Readable Time
	// ============================================================

	// time.Unix(seconds, nanoseconds)

	// Here:
	// unixTimestamp → seconds
	// 0 → nanoseconds

	readableTime := time.Unix(
		unixTimestamp,
		0,
	)

	fmt.Println(
		"Human Readable Time:",
		readableTime,
	)


	// ============================================================
	// 🔹 Formatting Readable Time
	// ============================================================

	// Format() converts time → formatted string

	fmt.Println(
		"Formatted Time:",
		readableTime.Format(
			"06-01-02 15:04:05",
		),
	)
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// Epoch Time / Unix Time
// → seconds since Jan 1 1970 UTC

// Unix()
// → time to unix timestamp

// time.Unix()
// → unix timestamp to time

// Format()
// → time to string

// Unix timestamps are useful because:
// - timezone independent
// - easy comparison
// - used in APIs/databases/logging