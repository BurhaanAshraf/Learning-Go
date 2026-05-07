package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
)

func cryptoHashing() {

	// ============================================================
	// 🔹 Hashing Basics
	// ============================================================

	// Hashing:
	// input → fixed-size irreversible output

	// SHA256 → 256-bit hash
	// SHA512 → 512-bit hash

	// Same input
	// → same hash

	// Tiny input change
	// → completely different hash
	// (Avalanche Effect)

	// Hashing is used for:
	// - passwords
	// - integrity checks
	// - verification


	// ============================================================
	// 🔹 SHA256 and SHA512
	// ============================================================

	password := "password1234"

	hashed256 := sha256.Sum256(
		[]byte(password),
	)

	hashed512 := sha512.Sum512(
		[]byte(password),
	)

	fmt.Println("Original Password:", password)

	// Raw binary hash bytes
	fmt.Println("SHA256 Raw:", hashed256)
	fmt.Println("SHA512 Raw:", hashed512)

	// %x converts bytes → readable hexadecimal
	fmt.Printf("SHA256 Hex: %x\n", hashed256)
	fmt.Printf("SHA512 Hex: %x\n", hashed512)


	// ============================================================
	// 🔹 Salting
	// ============================================================

	// Salt:
	// random value added before hashing

	// Same password + different salt
	// → different hash

	newPassword := "Password123"

	salt, err := generateSalt()

	if err != nil {
		panic("Cannot Generate Salt")
	}

	fmt.Printf("Salt Hex: %x\n", salt)

	signUpHash := hashPassword(
		newPassword,
		salt,
	)

	// Base64 used for safe storage/transmission
	saltString :=
		base64.StdEncoding.EncodeToString(salt)

	fmt.Println("Salt Base64:", saltString)

	fmt.Println("Signup Hash:", signUpHash)


	// ============================================================
	// 🔹 Login Verification
	// ============================================================

	decodedSalt, err :=
		base64.StdEncoding.DecodeString(
			saltString,
		)

	if err != nil {
		panic("Cannot Decode Salt")
	}

	loginHash := hashPassword(
		newPassword,
		decodedSalt,
	)

	if loginHash != signUpHash {
		panic("Invalid Credentials")
	}

	fmt.Println("Login Successful")
}


// ============================================================
// 🔹 Generate Salt
// ============================================================

func generateSalt() ([]byte, error) {

	salt := make([]byte, 16)

	// crypto/rand
	// → cryptographically secure randomness

	_, err := io.ReadFull(
		rand.Reader,
		salt,
	)

	if err != nil {
		return nil, err
	}

	return salt, nil
}


// ============================================================
// 🔹 Hash Password
// ============================================================

func hashPassword(
	password string,
	salt []byte,
) string {

	// Combine salt + password
	saltedPassword :=
		append(
			salt,
			[]byte(password)...,
		)

	hashedPassword :=
		sha512.Sum512(saltedPassword)

	// IMPORTANT:
	// hashedPassword is raw binary data
	// and not guaranteed to be valid UTF-8 text

	// so using:
	// string(hashedPassword[:])
	// may produce unreadable characters

	// Base64 converts binary data
	// into safe readable ASCII text

	// useful for:
	// - databases
	// - APIs
	// - JSON
	// - transmission
	// - logging

	return base64.StdEncoding.EncodeToString(
		hashedPassword[:],
	)
}


// ============================================================
// 🔹 IMPORTANT SECURITY NOTE
// ============================================================

// SHA256 / SHA512 alone
// are not ideal for password hashing

// Real systems use:
// - bcrypt
// - scrypt
// - argon2

// because they are intentionally slow


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// Hashing
// → fixed-size irreversible output

// SHA256
// → 256-bit hash

// SHA512
// → 512-bit hash

// Salt
// → random value before hashing

// crypto/rand
// → secure randomness

// Base64
// → binary → safe text

// string(hashBytes)
// → unsafe for raw binary hash data