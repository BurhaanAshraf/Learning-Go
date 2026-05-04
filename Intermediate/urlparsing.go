package main

import (
	"fmt"
	"net/url"
)

func URL() {

	// ============================================================
	// 🔹 URL Structure
	// ============================================================

	// [scheme://][userinfo@]host[:port][/path][?query][#fragment]

	rawURL := "https://example.com:8080/path?query=param#fragment"

	// Parse() → breaks URL into structured components
	parsedURL, err := url.Parse(rawURL)

	if err != nil {
		panic("Error parsing URL")
	}

	fmt.Println("Scheme:", parsedURL.Scheme)       // https
	fmt.Println("Host:", parsedURL.Host)           // example.com:8080
	fmt.Println("Port:", parsedURL.Port())         // 8080
	fmt.Println("Path:", parsedURL.Path)           // /path
	fmt.Println("Raw Query:", parsedURL.RawQuery)  // query=param
	fmt.Println("Fragment:", parsedURL.Fragment)   // fragment

	// Parsing = extracting structured information from URL


	// ============================================================
	// 🔹 Query Parameters (Reading)
	// ============================================================

	rawURL1 := "https://example.com/path?name=John&age=20"

	parsedURL1, err := url.Parse(rawURL1)

	if err != nil {
		panic("Error parsing URL")
	}

	// Query() → returns map[string][]string
	queryParams := parsedURL1.Query()

	fmt.Println("All Query Params:", queryParams)

	// Get() → fetch value by key
	fmt.Println("Name:", queryParams.Get("name"))
	fmt.Println("Age:", queryParams.Get("age"))


	// ============================================================
	// 🔹 Building URL (Using Query().Set())
	// ============================================================

	baseURL := &url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/path",
	}

	// Query() gives modifiable query map
	query := baseURL.Query()

	// Set() → adds OR replaces value
	query.Set("name", "John")
	query.Set("age", "20")

	// Encode() converts map → query string
	baseURL.RawQuery = query.Encode()

	fmt.Println("Built URL (Set):", baseURL.String())


	// ============================================================
	// 🔹 Building Query Using url.Values
	// ============================================================

	values := url.Values{}

	// Add() → appends values (does NOT replace)
	values.Add("name", "Jane")
	values.Add("age", "30")
	values.Add("city", "london")
	values.Add("country", "UK")

	// Encode() → converts to query string
	encodedQuery := values.Encode()

	fmt.Println("Encoded Query:", encodedQuery)

	// Manually building URL
	baseURL2 := "https://example.com/search"
	fullURL := baseURL2 + "?" + encodedQuery

	fmt.Println("Built URL (Values):", fullURL)
}


// ============================================================
// 🔹 DIFFERENCE: query.Set() vs values.Add()
// ============================================================

// query.Set(key, value)
// → replaces existing value if key already exists
// → used when modifying existing URL query

// Example:
// name=John → Set("name", "Jane")
// result → name=Jane



// values.Add(key, value)
// → appends new value
// → allows multiple values for same key

// Example:
// Add("name", "John")
// Add("name", "Jane")
// result → name=John&name=Jane



// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// url.Parse()
// → parse URL into components

// Query()
// → returns query params map

// Get()
// → fetch value by key

// Set()
// → add or replace query param

// Add()
// → add multiple values for same key

// Encode()
// → map → query string

// RawQuery
// → final query string in URL

// url.Values
// → helper for building query params