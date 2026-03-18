package main

import (
	"fmt"
	"maps"
)

func main() {

	// ============================================================
	// 🔹 1. Basics of Maps
	// ============================================================

	// - Maps store key-value pairs
	// - Keys must be UNIQUE
	// - Maps are UNORDERED (iteration order is not guaranteed)

	// Syntax:
	// var name map[keyType]valueType

	var table map[int]int
	_ = table


	// ============================================================
	// 🔹 2. Creating Maps
	// ============================================================

	// Using make()
	myMap := make(map[string]int)

	// Map literal (alternative)
	// myMap := map[string]int{
	//     "key1": 1,
	//     "key2": 2,
	// }

	// NOTE:
	// - make() initializes the map
	// - Without make → map is nil and cannot store values


	// ============================================================
	// 🔹 3. Adding & Updating Values
	// ============================================================

	myMap["key1"] = 9
	myMap["code"] = 18

	fmt.Println(myMap)
	fmt.Println(myMap["key1"])

	// Update existing key
	myMap["code"] = 30


	// ============================================================
	// 🔹 4. Accessing Values
	// ============================================================

	// If key does NOT exist → returns zero value of type
	// (int → 0, string → "", bool → false)

	fmt.Println(myMap["random"]) // 0


	// ============================================================
	// 🔹 5. Deleting Values
	// ============================================================

	delete(myMap, "key1")
	fmt.Println(myMap)

	// Add more values
	myMap["college"] = 96
	myMap["Address"] = 191201
	fmt.Println(myMap)

	// Clear entire map
	clear(myMap)
	fmt.Println(myMap)


	// ============================================================
	// 🔹 6. Checking Key Existence (IMPORTANT)
	// ============================================================

	myMap["key1"] = 10

	value, exists := myMap["key1"]
	fmt.Println(value, exists)

	// NOTE:
	// - exists → true if key present
	// - helps distinguish between "0 value" vs "not present"


	// ============================================================
	// 🔹 7. Map Literals & Comparison
	// ============================================================

	myMap1 := map[string]int{"a": 1, "b": 2, "c": 3}
	myMap2 := map[string]int{"a": 1, "b": 2, "c": 3}

	if maps.Equal(myMap1, myMap2) {
		fmt.Println("Equal")
	} else {
		fmt.Println("Not Equal")
	}

	// NOTE:
	// - maps.Equal compares maps (Go 1.21+)


	// ============================================================
	// 🔹 8. Iterating over Maps
	// ============================================================

	for k, v := range myMap1 {
		fmt.Printf("Key is %s and Value is %d\n", k, v)
	}

	// NOTE:
	// - Order is NOT guaranteed
	// - Use "_" if key or value not needed


	// ============================================================
	// 🔹 9. Nil Maps (IMPORTANT)
	// ============================================================

	var myNilMap map[string]string
	// nil map → cannot insert values (will panic)

	myNilMap = make(map[string]string) // now usable
	myNilMap["key"] = "hello"

	fmt.Println(myNilMap["key"])

	// NOTE:
	// - Reading from nil map → allowed (returns zero value)
	// - Writing → NOT allowed


	// ============================================================
	// 🔹 10. Length of Map
	// ============================================================

	fmt.Println(len(myMap1))


	// ============================================================
	// 🔹 11. Nested Maps (Map of Maps)
	// ============================================================

	OuterMaps := make(map[string]map[string]int)

	OuterMaps["Athish"] = myMap1
	OuterMaps["Burhaan"] = myMap

	fmt.Println(OuterMaps)

	// NOTE:
	// - Maps can hold other maps as values
	// - Useful for hierarchical data
}