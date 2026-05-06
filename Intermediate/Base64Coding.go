package main

import (
	"encoding/base64"
	"fmt"
)

func Base64Encoding() {

	// ============================================================
	// 🔹 Character Encoding Basics
	// ============================================================

	// ASCII
	// → 7-bit character encoding
	// → supports basic English characters

	// UTF-8
	// → variable-width encoding
	// → supports full Unicode
	// → backward compatible with ASCII
	// → most commonly used encoding today

	// UTF-16
	// → uses 16-bit units
	// → used internally in some systems


	// ============================================================
	// 🔹 What is Base64?
	// ============================================================

	// Base64 is:
	// binary → text encoding scheme

	// Converts binary data into readable ASCII text

	// Uses fixed set of 64 characters:
	// A-Z
	// a-z
	// 0-9
	// +
	// /
	// = (padding)

	// IMPORTANT:
	// Base64 is NOT encryption
	// because it can be decoded easily


	// ============================================================
	// 🔹 Why Base64 is Used
	// ============================================================

	// Used when binary data must travel through:
	// - text-based protocols
	// - URLs
	// - emails
	// - JSON/XML
	// - databases

	// Common uses:
	// - image embedding
	// - JWT tokens
	// - email attachments
	// - API data transfer


	// ============================================================
	// 🔹 URL Encoding
	// ============================================================

	// URL Encoding converts unsafe characters
	// into internet-safe representation

	// Example:
	// space → %20


	// ============================================================
	// 🔹 Original Binary Data
	// ============================================================

	// Strings internally become byte slices

	data := []byte("He~lo, Base64 Encoding")

	fmt.Println("Original Data:", data)
	fmt.Println("Original String:", string(data))


	// ============================================================
	// 🔹 Standard Base64 Encoding
	// ============================================================

	// EncodeToString()
	// → byte slice → base64 string

	encoded := base64.StdEncoding.EncodeToString(data)

	fmt.Println("Base64 Encoded:", encoded)


	// ============================================================
	// 🔹 Base64 Decoding
	// ============================================================

	// DecodeString()
	// → base64 string → byte slice

	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		panic("Error Decoding")
	}

	fmt.Println("Decoded Bytes:", decoded)

	// Convert byte slice → readable string
	fmt.Println("Decoded String:", string(decoded))


	// ============================================================
	// 🔹 URL Safe Base64 Encoding
	// ============================================================

	// Standard Base64 may contain:
	// +
	// /

	// These characters can create issues in URLs

	// URLEncoding replaces unsafe characters
	// to make output URL safe

	urlSafeEncoded :=
		base64.URLEncoding.EncodeToString(decoded)

	fmt.Println("URL Safe Encoding:", urlSafeEncoded)


	// ============================================================
	// 🔹 Padding
	// ============================================================

	// '=' is used as padding
	// to align encoded output properly

	// Important during decoding
	// because decoder expects proper structure


	// ============================================================
	// 🔹 Important Concepts
	// ============================================================

	// Base64:
	// → reversible encoding

	// NOT:
	// → encryption
	// → hashing

	// Use:
	// StdEncoding
	// → normal base64

	// Use:
	// URLEncoding
	// → URL-safe base64
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// ASCII
// → basic 7-bit encoding

// UTF-8
// → modern Unicode encoding

// Base64
// → binary to text encoding

// EncodeToString()
// → bytes → base64 string

// DecodeString()
// → base64 string → bytes

// StdEncoding
// → standard base64

// URLEncoding
// → URL-safe base64

// '='
// → padding character

// Base64 is NOT encryption
// → easily reversible