package main

import (
	"fmt"
	"time"
)

func Time() {

	// ============================================================
	// 🔹 Current Local Time
	// ============================================================

	// time.Now() gives current system time
	currentTime := time.Now()

	fmt.Println("Current Time:", currentTime)


	// ============================================================
	// 🔹 Creating Specific Time
	// ============================================================

	// time.Date() creates custom time manually

	specificTime := time.Date(
		2026,
		time.May,
		2,
		19,
		2,
		35,
		0,
		time.UTC,
	)

	fmt.Println("Specific Time:", specificTime)


	// ============================================================
	// 🔹 Parsing Time
	// ============================================================

	// Parse() converts string → time object

	// IMPORTANT:
	// Go uses reference date:
	// Mon Jan 02 15:04:05 MST 2006

	parsedDate1, _ := time.Parse(
		"2006-01-02",
		"2020-05-12",
	)

	parsedDate2, _ := time.Parse(
		"06-01-02",
		"26-04-02",
	)

	parsedDate3, _ := time.Parse(
		"06-1-2",
		"20-12-31",
	)

	fmt.Println("Parsed Date 1:", parsedDate1)
	fmt.Println("Parsed Date 2:", parsedDate2)
	fmt.Println("Parsed Date 3:", parsedDate3)


	// ============================================================
	// 🔹 Formatting Time
	// ============================================================

	// Format() converts time → string

	fmt.Println(
		"Formatted Current Time:",
		currentTime.Format("Mon 06-01-02 15-04-05"),
	)


	// ============================================================
	// 🔹 Adding Duration
	// ============================================================

	// Add() adds duration to time

	nextDay := currentTime.Add(time.Hour * 24)

	fmt.Println("One Day Later:", nextDay)

	// Weekday() gives day name
	fmt.Println("Weekday:", nextDay.Weekday())


	// ============================================================
	// 🔹 Time Zones
	// ============================================================

	// LoadLocation() loads timezone

	indiaLocation, _ := time.LoadLocation("Asia/Kolkata")

	utcTime := time.Date(
		2026,
		time.May,
		2,
		14,
		50,
		30,
		0,
		time.UTC,
	)

	// In() converts time to another timezone
	indiaTime := utcTime.In(indiaLocation)

	fmt.Println("UTC Time:", utcTime)
	fmt.Println("India Time:", indiaTime)


	// ============================================================
	// 🔹 Rounding + Truncating
	// ============================================================

	// Round() rounds to nearest duration
	roundedUTC := utcTime.Round(time.Hour)

	// Convert rounded time to local timezone
	roundedIndiaTime := roundedUTC.In(indiaLocation)

	fmt.Println("Rounded UTC Time:", roundedUTC)
	fmt.Println("Rounded India Time:", roundedIndiaTime)

	// Truncate() always cuts downward
	fmt.Println(
		"Truncated UTC Time:",
		utcTime.Truncate(time.Hour),
	)


	// ============================================================
	// 🔹 Another Time Zone
	// ============================================================

	newYorkLocation, _ := time.LoadLocation("America/New_York")

	newYorkTime := time.Now().In(newYorkLocation)

	fmt.Println("New York Time:", newYorkTime)


	// ============================================================
	// 🔹 Time Difference
	// ============================================================

	// Sub() returns difference between times

	minuteDifference := currentTime.Sub(
		currentTime.Add(-time.Minute * 26),
	)

	fmt.Println(
		"26 Minute Difference:",
		minuteDifference,
	)

	// Difference between two different times
	timeDifference := currentTime.Sub(indiaTime)

	fmt.Println(
		"Difference Between Current and India Time:",
		timeDifference,
	)


	// ============================================================
	// 🔹 Comparing Times
	// ============================================================

	// Comparison happens using actual instant in time
	// internally Go compares in UTC form

	fmt.Println(
		"Is currentTime after indiaTime?",
		currentTime.After(indiaTime),
	)

	// Before() checks opposite comparison
	fmt.Println(
		"Is currentTime before indiaTime?",
		currentTime.Before(indiaTime),
	)
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// time.Now()
// → current time

// time.Date()
// → create custom time

// time.Parse()
// → string to time

// Format()
// → time to string

// Add()
// → add duration

// Sub()
// → difference between times

// In()
// → convert timezone

// LoadLocation()
// → load timezone

// Round()
// → nearest duration

// Truncate()
// → cut downward

// After() / Before()
// → compare actual time instants