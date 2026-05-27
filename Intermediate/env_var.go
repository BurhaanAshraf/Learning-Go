package main

import (
	"fmt"
	"os"
	"strings"
)

func env_var() {

	// ============================================================
	// 🔹 Environment Variables
	// ============================================================

	// Environment Variables:
	// → system-wide key value pairs

	// Used for:
	// - configuration
	// - secrets
	// - API keys
	// - system paths
	// - runtime settings


	// ============================================================
	// 🔹 Getting Environment Variables
	// ============================================================

	// Getenv()
	// → fetches environment variable value

	userVar := os.Getenv("USER")

	homeVar := os.Getenv("HOME")

	fmt.Println("User Env Var:", userVar)

	fmt.Println("Home Env Var:", homeVar)

	// IMPORTANT:
	//
	// if key does not exist,
	// Getenv() returns:
	//
	// empty string ""


	// ============================================================
	// 🔹 Setting Environment Variables
	// ============================================================

	// Setenv()
	// → creates/updates environment variable

	key := "FRUIT"

	val := "APPLE"

	err := os.Setenv(key, val)

	if err != nil {
		panic(err)
	}

	fmt.Println("\nEnv var set on key FRUIT")

	fmt.Println("FRUIT Env Var:", os.Getenv("FRUIT"))


	// ============================================================
	// 🔹 Listing All Environment Variables
	// ============================================================

	// Environ()
	// → returns all environment variables

	// Format:
	//
	// KEY=VALUE

	for _, e := range os.Environ() {

		// SplitN(string, separator, n)

		// Split only at FIRST "="
		// because value may also contain "="

		kvPair := strings.SplitN(e, "=", 2)

		fmt.Println("Key:", kvPair[0])

		// safe check for value

		if len(kvPair) > 1 {

			fmt.Println("Value:", kvPair[1])
		}

		fmt.Println()
	}


	// ============================================================
	// 🔹 SplitN() Notes
	// ============================================================

	// n = 1
	// → no splitting

	// n = 2
	// → split at first separator

	// n = 3
	// → split at first two separators

	// n = -1
	// → split all separators

	// n = 0
	// → returns empty slice


	// ============================================================
	// 🔹 Removing Environment Variables
	// ============================================================

	// Unsetenv()
	// → removes environment variable

	err = os.Unsetenv("FRUIT")

	if err != nil {
		panic(err)
	}

	fmt.Println("Unset env var done on key FRUIT")

	fmt.Println("FRUIT Env Var:", os.Getenv("FRUIT"))

	fmt.Println("--------------------------------------------------")


	// ============================================================
	// 🔹 strings.SplitN()
	// ============================================================

	str := "a=b=c=d=e"

	fmt.Println(strings.SplitN(str, "=", -1))

	fmt.Println(strings.SplitN(str, "=", 0))

	fmt.Println(strings.SplitN(str, "=", 1))

	fmt.Println(strings.SplitN(str, "=", 2))

	fmt.Println(strings.SplitN(str, "=", 3))

	fmt.Println(strings.SplitN(str, "=", 5))
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// os.Getenv()
// → get env variable

// os.Setenv()
// → set env variable

// os.Unsetenv()
// → remove env variable

// os.Environ()
// → all environment variables

// Environment variables:
// → KEY=VALUE pairs

// strings.SplitN()
// → controlled splitting

// n = -1
// → split all

// n = 0
// → empty slice

// n = 1
// → no splitting

// n = 2
// → split once