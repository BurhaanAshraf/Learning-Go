package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func StringsAdv() {

	// ============================================================
	// 🔹 String Basics
	// ============================================================

	str := "Hello Go!"
	str2 := "I am Burhaan"

	// String concatenation joins strings together
	str3 := str + " " + str2

	// %c prints character at that byte position
	fmt.Printf("%c\n", str3[0])

	// Direct indexing gives byte / ASCII value
	fmt.Println(str[0])

	// String slicing gives substring
	fmt.Println(str[1 : len(str)-1])


	// ============================================================
	// 🔹 String Conversion
	// ============================================================

	num := 18

	// strconv.Itoa converts int → string
	str4 := strconv.Itoa(num)

	fmt.Println(len(str4))


	// ============================================================
	// 🔹 Split
	// ============================================================

	fruits := "apple , banana , cherry , pear"
	fruits2 := "apple-banana-cherry-pear"

	// Split breaks one string into slice of strings
	parts := strings.Split(fruits, ",")
	parts1 := strings.Split(fruits2, "-")

	fmt.Println(parts)
	fmt.Println(parts1)


	// ============================================================
	// 🔹 Join
	// ============================================================

	// Join converts slice → single string

	countries := []string{
		"germany",
		"italy",
		"swiss",
		"france",
	}

	joinedCountries := strings.Join(countries, ",")

	fmt.Println(joinedCountries)


	// ============================================================
	// 🔹 Contains + Replace
	// ============================================================

	// Contains checks if substring exists
	fmt.Println(
		strings.Contains(joinedCountries, "germany"),
	)

	// Replace changes old substring to new substring
	replaced := strings.Replace(
		joinedCountries,
		"germany",
		"India",
		1,
	)

	fmt.Println(replaced)


	// ============================================================
	// 🔹 TrimSpace
	// ============================================================

	str5 := "     Hello Everyone!"

	fmt.Println(str5)

	// TrimSpace removes leading + trailing spaces
	str5 = strings.TrimSpace(str5)

	fmt.Println(str5)


	// ============================================================
	// 🔹 Uppercase + Lowercase
	// ============================================================

	// Converts entire string to uppercase
	fmt.Println(strings.ToUpper(str5))

	str6 := "HELLO EVERYONE!"

	// Converts entire string to lowercase
	str6 = strings.ToLower(str6)

	fmt.Println(str6)


	// ============================================================
	// 🔹 Repeat
	// ============================================================

	// Repeat repeats string n times
	fmt.Println(strings.Repeat("foo ", 3))


	// ============================================================
	// 🔹 Count
	// ============================================================

	// Count counts occurrences of substring
	fmt.Println(
		strings.Count("India", "Ind"),
	)


	// ============================================================
	// 🔹 Prefix + Suffix
	// ============================================================

	// HasPrefix checks starting characters
	fmt.Println(
		strings.HasPrefix("Hello", "He"),
	)

	// HasSuffix checks ending characters
	fmt.Println(
		strings.HasSuffix("Hello", "l"),
	)


	// ============================================================
	// 🔹 Regular Expressions
	// ============================================================

	newStr := "1234567 Hello World Golang 890"

	// \d+ means one or more digits
	re := regexp.MustCompile(`\d+`)

	// Finds all number matches
	matches := re.FindAllString(newStr, -1)

	fmt.Println(matches)


	// ============================================================
	// 🔹 Rune Count
	// ============================================================

	// len() counts bytes
	// RuneCount counts actual characters

	fmt.Println(
		utf8.RuneCountInString(newStr),
	)


	// ============================================================
	// 🔹 String Builder
	// ============================================================

	// strings.Builder is efficient for building strings
	// better than using + many times

	var builder strings.Builder

	builder.WriteString("Hello")
	builder.WriteRune(' ')
	builder.WriteString("World")

	result := builder.String()

	fmt.Println(result)

	// Reset clears old content
	builder.Reset()

	builder.WriteString(
		"This is after resetting the builder",
	)

	result = builder.String()

	fmt.Println(result)
}


// ============================================================
// 🔹 Quick Revision
// ============================================================

// len() → byte length

// RuneCountInString() → actual character count

// str[i] → byte value

// %c → prints character

// str[a:b] → substring

// strconv.Itoa() → int to string

// Split() → string to slice

// Join() → slice to string

// Contains() → substring check

// Replace() → replace substring

// TrimSpace() → remove spaces

// ToUpper() / ToLower()

// Repeat() → repeat string

// Count() → count substring

// HasPrefix() / HasSuffix()

// regexp → pattern matching

// strings.Builder → efficient string creation