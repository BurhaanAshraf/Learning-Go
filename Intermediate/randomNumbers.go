package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// ============================================================
// 🔹 Dice Function
// ============================================================

// Intn(6) gives:
// [0, 6)
//
// so +1 makes range:
// [1, 6]

func rollDice() int {
	return rand.Intn(6) + 1
}


func mathRand() {

	// ============================================================
	// 🔹 Random Numbers
	// ============================================================

	// Intn(n)
	// generates random number in:
	// [0, n)
	//
	// called half-open interval

	randomNum := rand.Intn(101)

	fmt.Println(
		"Random Number [0-100]:",
		randomNum,
	)

	// [1,100]
	randomNum2 := rand.Intn(100) + 1

	fmt.Println(
		"Random Number [1-100]:",
		randomNum2,
	)


	// ============================================================
	// 🔹 Fixed Seed
	// ============================================================

	// Seed fixes pseudo-random sequence

	// Same seed
	// → same output every time

	fixedSeed := rand.New(
		rand.NewSource(42),
	)

	fmt.Println(
		"Fixed Seed Random:",
		fixedSeed.Intn(101),
	)


	// ============================================================
	// 🔹 Dynamic Seed
	// ============================================================

	// Unix time changes every second
	// so random sequence changes

	dynamicSeed := rand.New(
		rand.NewSource(time.Now().Unix()),
	)

	fmt.Println(
		"Dynamic Seed Random:",
		dynamicSeed.Intn(100),
	)


	// ============================================================
	// 🔹 Random Float
	// ============================================================

	// Float64()
	// generates value in:
	// [0.0, 1)

	randomFloat := rand.Float64()

	fmt.Println(
		"Random Float:",
		randomFloat,
	)


	// ============================================================
	// 🔹 Dice Game
	// ============================================================

	for {

		fmt.Println("\n1. Roll Dice")
		fmt.Println("2. Exit")
		fmt.Println("Choose option:")

		var userChoice int

		_, err := fmt.Scan(&userChoice)

		if err != nil ||
			(userChoice != 1 && userChoice != 2) {

			fmt.Println("Please Enter Valid Option")
			continue
		}


		// ========================================================
		// 🔹 Exit Game
		// ========================================================

		if userChoice == 2 {
			fmt.Println("Thanks For Playing!")
			break
		}


		// ========================================================
		// 🔹 Roll Dice
		// ========================================================

		dice1 := rollDice()
		dice2 := rollDice()

		fmt.Println("Dice 1:", dice1)
		fmt.Println("Dice 2:", dice2)

		fmt.Println(
			"Sum of Dices:",
			dice1+dice2,
		)


		// ========================================================
		// 🔹 Play Again
		// ========================================================

		fmt.Println("Roll Again? (Y/N)")

		var playAgain string

		_, newErr := fmt.Scan(&playAgain)

		playAgain = strings.ToUpper(playAgain)

		if newErr != nil ||
			(playAgain != "Y" &&
				playAgain != "N" &&
				playAgain != "YES" &&
				playAgain != "NO") {

			fmt.Println("Invalid Input! Exiting...")
			break
		}

		if playAgain == "N" ||
			playAgain == "NO" {

			fmt.Println("Thanks For Playing!")
			break
		}
	}
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// rand.Intn(n)
// → random int in [0,n)

// +1
// → shift range upward

// Float64()
// → random float in [0.0,1)

// Seed
// → source for pseudo-random generator

// Same seed
// → same output

// time.Now().Unix()
// → changing seed using epoch time

// strings.ToUpper()
// → normalize user input

// ============================================================
// 🔹 math/rand vs crypto/rand
// ============================================================

// math/rand
// → pseudo-random generator
// → deterministic if seed is known
// → same seed gives same output

// Used for:
// - games
// - simulations
// - testing
// - random ordering

// NOT secure for authentication/security



// crypto/rand
// → cryptographically secure random generator
// → unpredictable output
// → used for security-sensitive systems

// Used for:
// - passwords
// - OTPs
// - API keys
// - authentication tokens
// - encryption



// IMPORTANT:

// math/rand
// → fast but predictable

// crypto/rand
// → secure but slower