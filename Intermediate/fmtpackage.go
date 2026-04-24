package main

import "fmt"

func FMT() {

	// -----------------------------
	// Formatting Functions
	// -----------------------------

	// Sprint → joins values into a string (no spaces added)
	a := fmt.Sprint("Hello", "World", 123, true)

	// Sprintln → joins values with spaces + adds newline
	b := fmt.Sprintln("Hello", "World", 123, true)

	fmt.Print(a)
	fmt.Println(b)

	// Sprintf → formats and returns a string
	age := 10
	name := "Burhaan"

	// %s = string, %d = integer
	c := fmt.Sprintf("Name is %s and Age is %d\n", name, age)
	fmt.Print(c)

	// -----------------------------
	// Scanning Functions
	// -----------------------------

	// Scan Example
	var scanName string
	var scanAge int

	// Scan → reads input separated by spaces
	// Example input: Burhaan 21
	fmt.Scan(&scanName, &scanAge)
	fmt.Printf("Using Scan -> Name: %s and Age: %d\n", scanName, scanAge)

	// Scanln Example
	var scanlnName string
	var scanlnAge int

	// Scanln → reads input until newline
	// Example input: Burhaan 22
	fmt.Scanln(&scanlnName, &scanlnAge)
	fmt.Printf("Using Scanln -> Name: %s and Age: %d\n", scanlnName, scanlnAge)

	// Scanf Example
	var scanfName string
	var scanfAge int

	// Scanf → reads input using format specifiers
	// Example input: Burhaan 23
	fmt.Scanf("%s %d", &scanfName, &scanfAge)
	fmt.Printf("Using Scanf -> Name: %s and Age: %d\n", scanfName, scanfAge)

	// -----------------------------
	// Error Formatting
	// -----------------------------

	var driveAge int

	// Input age for driving check
	fmt.Scan(&driveAge)

	// Function may return an error
	err := checkDrive(driveAge)

	// nil means no error
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// checkDrive → validates driving age
func checkDrive(age int) error {

	// Error if age is below 18
	if age < 18 {
		// Errorf → creates formatted error
		return fmt.Errorf("Age %d is too young to drive", age)
	}

	return nil
}