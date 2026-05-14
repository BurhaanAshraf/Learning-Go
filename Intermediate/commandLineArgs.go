package main

import (
	"flag"
	"fmt"
	"os"
)

func CLA() {

	// ============================================================
	// 🔹 Command Line Arguments
	// ============================================================

	// Command-line arguments:
	// → values passed from terminal while running program

	// Example:
	//
	// go run main.go -name Burhaan -age 20 -male=false

	// Internally:
	// arguments are stored inside:
	//
	// os.Args


	// ============================================================
	// 🔹 os.Args
	// ============================================================

	// os.Args[0]
	// → path of currently running executable

	// IMPORTANT:
	//
	// when using:
	// go run main.go
	//
	// Go first creates temporary executable
	// inside temporary build directory

	// Example:
	// /tmp/go-build123/b001/exe/main

	// Then Go runs executable

	// After execution:
	// temporary build files are auto deleted

	fmt.Println("Executable Path:", os.Args[0])


	// ============================================================
	// 🔹 Declaring Variables
	// ============================================================

	var userName string
	var userAge int
	var isMale bool


	// ============================================================
	// 🔹 String Flags
	// ============================================================

	// StringVar(
	//     variableAddress,
	//     flagName,
	//     defaultValue,
	//     description,
	// )

	flag.StringVar(
		&userName,
		"name",
		"John",
		"name of the user",
	)

	// &userName
	// → address of variable
	// → flag package directly updates variable


	// ============================================================
	// 🔹 Integer Flags
	// ============================================================

	flag.IntVar(
		&userAge,
		"age",
		18,
		"age of the user",
	)


	// ============================================================
	// 🔹 Boolean Flags
	// ============================================================

	flag.BoolVar(
		&isMale,
		"male",
		true,
		"gender of user",
	)

	// Writing:
	//
	// -male false
	//
	// does NOT work correctly because:
	//
	// -male alone already means true
	//
	// and "false" becomes separate argument


	// ============================================================
	// 🔹 Parsing Flags
	// ============================================================

	// Parse():
	// → reads terminal arguments
	// → assigns values to variables

	flag.Parse()


	// ============================================================
	// 🔹 Printing Parsed Values
	// ============================================================

	fmt.Println("\nParsed Values:")

	fmt.Println("Name:", userName)
	fmt.Println("Age:", userAge)
	fmt.Println("Male:", isMale)

// ============================================================
// 🔹 Remaining Non-Flag Arguments
// ============================================================

// Args()
// → returns leftover arguments
// after flag parsing

// Example:
//
// go run main.go -name Burhaan extra1 extra2
//
// extra1 and extra2 become:
//
// flag.Args()

remainingArgs := flag.Args()

fmt.Println(
	"Remaining Arguments:",
	remainingArgs,
)
}


// ============================================================
// 🔹 HOW TO RUN
// ============================================================

// Using default values:
//
// go run main.go


// Custom values:
//
// go run main.go -name Burhaan -age 20 -male=false


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// os.Args
// → command-line arguments

// os.Args[0]
// → executable path

// go run
// → creates temporary executable
// → runs it
// → auto deletes temp build files

// flag package
// → command-line flag parsing

// StringVar()
// → string flag

// IntVar()
// → integer flag

// BoolVar()
// → boolean flag

// flag.Parse()
// → parse terminal arguments

// Boolean flags:
// presence itself means true

// &variable
// → variable address (pointer)