package main

import "fmt"

func Arrays() {
	

	// ===============================
	// 1. Array Declaration
	// ===============================
	// Syntax: var arrayName [size]elementType
	// When an array is declared without values, Go initializes it with zero values

	var numbers [5]int
	numbers[4] = 5 // update value at index 4

	fmt.Println("Numbers array:", numbers)

	// ===============================
	// 2. Array Declaration + Initialization
	// ===============================

	var names [4]string = [4]string{"Burhaan", "Farhaan", "Rizwaan", "Epstein"}
	fmt.Println("Names array:", names)

	// ===============================
	// 3. Short Declaration (:=)
	// ===============================
	// Go infers the type automatically

	newNames := [4]string{"Adam", "Joseph", "Raven", "Claw"}
	fmt.Println("New Names:", newNames)

	// ===============================
	// 4. Array Copy Behaviour
	// ===============================
	// Arrays in Go are VALUE TYPES
	// When assigned, the entire array is copied

	originalArray := [4]int{1, 2, 3, 4}
	copiedArray := originalArray

	// modifying original array does NOT affect copied array
	originalArray[0] = 100

	fmt.Println("Original Array:", originalArray)
	fmt.Println("Copied Array:", copiedArray)

	// ===============================
	// 5. Iterating using classic for loop
	// ===============================

	for i := 0; i < len(numbers); i++ {
		fmt.Printf("Index: %d , Value: %d\n", i, numbers[i])
	}

	// ===============================
	// 6. Iterating using range
	// ===============================
	// range returns two values:
	// index and value

	for index, value := range copiedArray {
		fmt.Printf("Index: %d , Value: %d\n", index, value)
	}

	// ===============================
	// 7. Blank Identifier (_)
	// ===============================
	// Used when we want to ignore a returned value

	for _, value := range numbers {
		fmt.Printf("Value: %d\n", value)
	}

	// Go does not allow unused variables
	b := 2

	// Assigning to blank identifier avoids the error
	_ = b

	// ===============================
	// 8. Length of an Array
	// ===============================

	fmt.Println("Length of numbers array:", len(numbers))

	// ===============================
	// 9. Comparing Arrays
	// ===============================
	// Arrays can be compared if they have
	// same type and same size

	array1 := [3]int{1, 2, 3}
	array2 := [3]int{1, 2, 3}

	if array1 == array2 {
		fmt.Println("Array1 and Array2 are equal")
	} else {
		fmt.Println("Array1 and Array2 are not equal")
	}

	// ===============================
	// 10. Multidimensional Arrays
	// ===============================

	var matrix [3][3]int = [3][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	fmt.Println("Matrix:", matrix)

	// ===============================
	// 11. Pointer to Array
	// ===============================
	// A pointer can store the address of an array
	//A pointer is a variable that stores the address of another variable.

	newOriginalArray := [3]int{1, 2, 3}

	var newCopiedArray *[3]int
	newCopiedArray = &newOriginalArray

	// Alternative ways to declare pointer
	// newCopiedArray := &newOriginalArray
	// var newCopiedArray = &newOriginalArray

	fmt.Println(newOriginalArray)
	fmt.Println(*newCopiedArray)

	// modifying via pointer changes original array
	newCopiedArray[0] = 10
	newOriginalArray[1] = 100

	fmt.Println("Original Array via pointer:", newOriginalArray)
	fmt.Println("Dereferenced Pointer:", *newCopiedArray)
}