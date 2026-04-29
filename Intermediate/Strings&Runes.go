package main

import (
	"fmt"
	"unicode/utf8"
)

func Strings() {

	message := "Hello Go"
	_ = message

	rawMessage := `Hello \nGo`
	_ = rawMessage

	name := ""

	greetings := "Good Morning, " + name
	_ = greetings

	str1 := "apple"
	str2 := "banana"

	fmt.Println(str1 < str2)

	str3 := "app"
	fmt.Println(str3 < str1)

	lowercaseApple := "apple"
	uppercaseApple := "APPLE"

	fmt.Println(lowercaseApple > uppercaseApple)

	// ============================================================
	// 🔹 String Iteration
	// ============================================================

	for i, char := range message {
		fmt.Printf("Index is %d and Hex Value is %x\n", i, char)
	}

	fmt.Println(utf8.RuneCountInString(lowercaseApple))


	// ============================================================
	// 🔹 Runes
	// ============================================================

	var ch rune = 'a'
	_ = ch

	jch := '日'

	fmt.Println(jch)
	fmt.Printf("%c\n", jch)

	convertedStr := string(ch)
	fmt.Println(convertedStr)
	fmt.Printf("%T\n", convertedStr)

	const NIHONGO = "日本語"
	fmt.Println(NIHONGO)

	const JHELLO = "こんにちは"

	for _, char := range JHELLO {
		fmt.Printf("%c\n", char)
	}

	smiley := '😊'
	fmt.Println(smiley)
	fmt.Printf("%c\n", smiley)
}


// ============================================================
// 🔹 String Basics
// ============================================================

// A string is a sequence of bytes (UTF-8 encoded).
// Each byte is of type uint8.

// Strings are immutable → cannot be changed after creation.


// ============================================================
// 🔹 Raw Strings
// ============================================================

// Backticks (`) create raw string literals.
// Escape sequences like \n are NOT processed.


// ============================================================
// 🔹 String Length
// ============================================================

// len(str) → returns number of BYTES, not characters.

// Example:
// "Hello" → 5 bytes
// "日本語" → more than 3 bytes (multi-byte characters)


// ============================================================
// 🔹 Indexing
// ============================================================

// str[i] → returns byte value (uint8), NOT character.

// Example:
// rawMessage[0] → ASCII value of first byte


// ============================================================
// 🔹 String Comparison
// ============================================================

// Strings are compared lexicographically (dictionary order).
// Comparison is based on Unicode/ASCII values.

// Example:
// "apple" < "banana" → true
// lowercase > uppercase (ASCII values)


// ============================================================
// 🔹 String Concatenation
// ============================================================

// Use + operator to join strings.

// Note:
// Go does NOT add space automatically.


// ============================================================
// 🔹 String Iteration
// ============================================================

// range over string → gives index + rune (Unicode character).

// index → byte position
// value → rune (int32)


// ============================================================
// 🔹 Rune (VERY IMPORTANT)
// ============================================================

// rune = alias for int32
// represents a Unicode code point (character).

// Go does NOT have char type → uses rune instead.


// ============================================================
// 🔹 Rune vs Byte
// ============================================================

// byte → 1 byte (uint8)
// rune → 4 bytes (int32)

// ASCII characters → 1 byte
// Unicode characters → multiple bytes


// ============================================================
// 🔹 Counting Characters
// ============================================================

// utf8.RuneCountInString(str)
// returns number of characters (runes), NOT bytes


// ============================================================
// 🔹 Conversion
// ============================================================

// string(rune) → converts rune to string


// ============================================================
// 🔹 Unicode Support
// ============================================================

// Go fully supports Unicode:
// Multilingual, emojis, symbols, etc.


// ============================================================
// 🔹 Important Notes
// ============================================================

// Strings are immutable
// len() gives bytes, not characters
// range gives runes (correct way to iterate characters)
// indexing gives bytes