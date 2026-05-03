package main

import (
	"fmt"
	"time"
)

func Parse() {

	// ============================================================
	// 🔹 Time Parsing
	// ============================================================

	// Parse() converts string → time object

	// IMPORTANT:
	// Go uses reference time:
	//
	// Mon Jan 02 15:04:05 MST 2006
	//
	// These exact numbers must stay same:
	// 2006 → year
	// 01   → month
	// 02   → day
	// 15   → 24hr hour
	// 03   → 12hr hour
	// 04   → minute
	// 05   → second


	// ============================================================
	// 🔹 Layout 1 (Timezone Offset Format)
	// ============================================================

	// -07:00 represents timezone offset
	// +05:30 means UTC + 5hr 30min

	layout1 := "2006-01-02 15:04:05 -07:00"

	str1 := "2024-07-04 14:30:18 +05:30"

	parsedTime1, err := time.Parse(layout1, str1)

	if err != nil {
		panic(err)
	}

	fmt.Println(
		"Parsed Time with Offset Layout:",
		parsedTime1,
	)


	// ============================================================
	// 🔹 Layout 2 (RFC3339 Format)
	// ============================================================

	// T separates date and time
	// Z means UTC timezone

	layout2 := "2006-01-02T15:04:05Z07:00"

	str2 := "2024-07-04T14:30:18Z"

	parsedTime2, err := time.Parse(layout2, str2)

	if err != nil {
		panic(err)
	}

	fmt.Println(
		"Parsed RFC3339 Time:",
		parsedTime2,
	)


	// ============================================================
	// 🔹 Custom Human Readable Format
	// ============================================================

	// Jan → month name
	// 02 → day
	// 2006 → year
	// 03 → 12hr hour
	// PM → AM/PM indicator

	layout3 := "Jan 02, 2006 03:04 PM"

	str3 := "May 03, 2026 11:10 AM"

	parsedTime3, err := time.Parse(layout3, str3)

	if err != nil {
		panic(err)
	}

	fmt.Println(
		"Parsed Custom Time:",
		parsedTime3,
	)
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// time.Parse()
// → string to time object

// Layout must match string exactly

// Go uses reference date:
// Mon Jan 02 15:04:05 MST 2006

// 15 → 24hr format

// 03 PM → 12hr format

// -07:00 / Z07:00
// → timezone offset

// T
// → separates date and time

// Z
// → UTC timezone