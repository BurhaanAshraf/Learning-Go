package main

import (
	"flag"
	"fmt"
	"os"
)

func commandLineSubArguments() {

	// ============================================================
	// 🔹 Simple Flag Using flag.String()
	// ============================================================

	// String()
	// → creates variable internally
	// → returns POINTER to string

	// Syntax:
	//
	// flag.String(
	//     flagName,
	//     defaultValue,
	//     description,
	// )

	newStringFlag := flag.String("guest", "Kartik", "Juniors Ka Bade Bhaiya")

	// Parse command-line flags
	flag.Parse()

	// IMPORTANT:
	// String() returns *string
	// so we use:
	// *newStringFlag
	// to access actual value

	fmt.Println("Guest Flag:", *newStringFlag)


	// ============================================================
	// 🔹 StringVar()
	// ============================================================

	// StringVar()
	// → directly stores parsed value
	// inside existing variable

	var stringFlag string

	flag.StringVar(&stringFlag, "user", "John", "Name of User")

	// IMPORTANT:
	//
	// &stringFlag
	// → address of variable
	//
	// flag package needs address
	// so it can MODIFY original variable


	// ============================================================
	// 🔹 Subcommands
	// ============================================================

	// NewFlagSet()
	// → independent command parser

	// Similar to:
	//
	// git push
	// git commit
	//
	// where:
	// push and commit are subcommands

	subcommand1 := flag.NewFlagSet("firstsub", flag.ExitOnError)

	subcommand2 := flag.NewFlagSet("secondsub", flag.ExitOnError)


	// ============================================================
	// 🔹 Flags For firstsub
	// ============================================================

	firstFlag := subcommand1.Bool("processing", false, "Command Processing Status")

	secondFlag := subcommand1.Int("bytes", 1024, "Byte Length of Result")


	// ============================================================
	// 🔹 Flags For secondsub
	// ============================================================

	flagSC2 := subcommand2.String("language", "Go", "Enter your language")


	// ============================================================
	// 🔹 Argument Validation
	// ============================================================

	// os.Args[1]
	// → first argument after executable name

	if len(os.Args) < 2 {

		fmt.Println("This program requires subcommands")

		os.Exit(1)
	}


	// ============================================================
	// 🔹 Subcommand Routing
	// ============================================================

	switch os.Args[1] {


	// ========================================================
	// 🔹 firstsub
	// ========================================================

	case "firstsub":

		// Parse only arguments
		// belonging to firstsub

		subcommand1.Parse(os.Args[2:])

		fmt.Println("\nSubCommand1:")

		// IMPORTANT:
		// Bool()/Int()/String()
		// return pointers

		fmt.Println("Processing:", *firstFlag)

		fmt.Println("Bytes:", *secondFlag)


	// ========================================================
	// 🔹 secondsub
	// ========================================================

	case "secondsub":

		subcommand2.Parse(os.Args[2:])

		fmt.Println("\nSubCommand2:")

		fmt.Println("Language:", *flagSC2)


	// ========================================================
	// 🔹 Invalid Command
	// ========================================================

	default:

		fmt.Println("NO VALID SUBCOMMAND ENTERED")

		os.Exit(1)
	}
}


// ============================================================
// 🔹 HOW TO RUN
// ============================================================

// Simple flag:
//
// go run main.go -guest Burhaan


// First subcommand:
//
// go run main.go firstsub -processing=true -bytes=2048


// Second subcommand:
//
// go run main.go secondsub -language Python


// ============================================================
// 🔹 POINTER NOTES FOR FLAGS
// ============================================================

// StringVar():
// → needs ADDRESS of variable

// because flag package must MODIFY
// original variable internally

// Example:
//
// var userName string
//
// flag.StringVar(&userName, ...)


// &userName
// → address of variable


// ============================================================
// 🔹 flag.String()
// ============================================================

// flag.String()
// → creates variable internally
// → returns POINTER to variable

// Example:
//
// newFlag := flag.String(...)


// Type:
//
// *string


// ============================================================
// 🔹 Dereferencing
// ============================================================

// Since flag.String()
// returns pointer,
// we use:
//
// *newFlag
//
// to access actual value


// *pointer
// → actual value stored at address


// ============================================================
// 🔹 Easy Mental Model
// ============================================================

// &variable
// → give address

// *pointer
// → get value from address


// ============================================================
// 🔹 Difference Summary
// ============================================================

// StringVar():
//
// YOU create variable
// flag updates it using address


// String():
//
// flag package creates variable internally
// and returns pointer


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// flag.String()
// → returns pointer

// flag.StringVar()
// → stores directly in variable

// flag.Parse()
// → parses flags

// NewFlagSet()
// → independent subcommand parser

// os.Args
// → command-line arguments

// os.Args[1]
// → first subcommand

// Bool()/Int()/String()
// → return pointers

// *pointer
// → actual value

// &variable
// → address of variable

// os.Exit(1)
// → terminate program with error