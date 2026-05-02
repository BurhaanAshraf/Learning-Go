package main

import (
	"fmt"
	"regexp"
)

func regex() {

	// ============================================================
	// 🔹 REGEX (Regular Expressions)
	// ============================================================

	// Regex is used for:
	// - validation
	// - searching
	// - pattern matching
	// - replacing text


	// ============================================================
	// 🔹 [] vs ()
	// ============================================================

	// [] → Character Set
	// Matches ONE character from given options

	// Examples:
	// [aeiou] → one vowel
	// [a-z] → one lowercase letter
	// [0-9] → one digit


	// () → Capturing Group
	// Groups pattern together
	// and stores matched part separately

	// Capturing means:
	// match + save specific parts for later use

	// Example:
	// (\d{2})
	// captures 2 digits separately


	// ============================================================
	// 🔹 Email Validation
	// ============================================================

	// [a-zA-Z0-9._+%-]+ → username
	// @ → required symbol
	// [a-zA-Z0-9.-]+ → domain
	// \. → actual dot
	// [a-zA-Z]{2,} → extension

	re := regexp.MustCompile(
		`[a-zA-Z0-9._+%-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`,
	)

	email1 := "user@email.com"
	email2 := "Invalid Email"

	// MatchString() → checks pattern match
	fmt.Println(re.MatchString(email1))
	fmt.Println(re.MatchString(email2))


	// ============================================================
	// 🔹 Capturing Groups
	// ============================================================

	// (\d{2}) → captures day
	// (\d{2}) → captures month
	// (\d{4}) → captures year

	re = regexp.MustCompile(`(\d{2})-(\d{2})-(\d{4})`)

	date := "19-04-2006"

	// FindStringSubmatch()
	// returns:
	// [fullMatch, group1, group2, group3]

	submatches := re.FindStringSubmatch(date)

	fmt.Println(submatches)

	// Full match
	fmt.Println(submatches[0])

	// Captured groups
	fmt.Println(submatches[1]) // day
	fmt.Println(submatches[2]) // month
	fmt.Println(submatches[3]) // year


	// ============================================================
	// 🔹 Replace Using Regex
	// ============================================================

	str := "Hello World"

	// [aeiou] → matches vowels
	re = regexp.MustCompile(`[aeiou]`)

	// Replace vowels with *
	result := re.ReplaceAllString(str, "*")

	fmt.Println(result)


	// ============================================================
	// 🔹 Regex Flags
	// ============================================================

	// (?i) → case insensitive
	// (?m) → multiline mode
	// (?s) → dot matches newline

	re = regexp.MustCompile(`(?i)go`)

	str = "Hello Golang"

	fmt.Println(re.MatchString(str))
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// regexp.MustCompile()
// → compile regex

// MatchString()
// → check pattern match

// FindStringSubmatch()
// → get captured groups

// ReplaceAllString()
// → replace matches

// [] → one character from set

// () → capture/group pattern

// Capturing
// → match + store separately

// \d → digit

// + → one or more

// {n} → exact count

// (?i) → case insensitive