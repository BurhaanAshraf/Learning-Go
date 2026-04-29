package main

import "fmt"

func StringFormatting() {

	// ============================================================
	// 🔹 Number Formatting
	// ============================================================

	num := 42

	// %05d
	// d → integer
	// 5 → minimum width
	// 0 → fill remaining space with leading zeros

	// Output: 00042

	fmt.Printf("%05d\n", num)


	// ============================================================
	// 🔹 String Width Formatting
	// ============================================================

	message := "Hello"

	// %10s
	// s → string
	// 10 → minimum width

	// If string length < 10
	// leading spaces are added

	fmt.Printf("|%10s|\n", message)


	// %-10s
	// - means left aligned

	// If string length < 10
	// trailing spaces are added

	fmt.Printf("|%-10s|\n", message)


	// ============================================================
	// 🔹 Normal String vs Raw String
	// ============================================================

	// Normal string uses escape sequences

	message2 := "Hello \nWorld"

	// Raw string uses backticks `
	// Escape sequences are treated as plain text

	message3 := `Hello \nWorld`

	fmt.Println(message2)
	fmt.Println(message3)
}


// ============================================================
// 🔹 Quick Revision
// ============================================================

// %d → integer

// %s → string

// %05d → pad with leading zeros

// %10s → right aligned string

// %-10s → left aligned string

// "" → normal string (escape sequences work)

// `` → raw string (literal text)