package main

import "fmt"

func recursion() {

	fmt.Println(factorial(5))
	fmt.Println(factorial(10))

	fmt.Println(sumofdigits(291))
}

func factorial(n int) int {

	// Base Case
	// Stops recursion when n becomes 0
	if n == 0 {
		return 1
	}

	// Recursive Case
	// Function calls itself with smaller input
	return n * factorial(n-1)
}

func sumofdigits(num int) int {

	// Base Case
	// When number becomes single digit
	if num < 10 {
		return num
	}

	// Recursive Case
	// Breaks number into last digit + remaining digits
	return num%10 + sumofdigits(num/10)
}


// ============================================================
// Recursion Definition
// ============================================================

// Recursion is a programming technique where a function calls
// itself to solve smaller instances of the same problem.


// ============================================================
// Key Components of Recursion
// ============================================================

// 1. Base Case
// The condition that stops the recursion.

// 2. Recursive Case
// The part where the function calls itself with smaller input.


// ============================================================
// Example Flow: factorial(5)
// ============================================================

// factorial(5)
// = 5 * factorial(4)
// = 5 * 4 * factorial(3)
// = 5 * 4 * 3 * factorial(2)
// = 5 * 4 * 3 * 2 * factorial(1)
// = 5 * 4 * 3 * 2 * 1 * factorial(0)
// = 120


// ============================================================
// Example Flow: sumofdigits(291)
// ============================================================

// sumofdigits(291)
// = 1 + sumofdigits(29)
// = 1 + 9 + sumofdigits(2)
// = 1 + 9 + 2
// = 12


// ============================================================
// Important Notes
// ============================================================

// Every recursive function must have a base case.
// Without a base case recursion continues indefinitely,
// leading to stack overflow.

// Recursive calls must move toward the base case.


// ============================================================
// When Recursion is Useful
// ============================================================

// Tree traversal
// Divide and conquer algorithms
// Backtracking
// Mathematical problems (factorial, Fibonacci)