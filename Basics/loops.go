package main

import (
	"fmt"
)
func loops() {

	for i:= 1 ; i <= 5 ; i++ { // this is simple iteration over range
		fmt.Println(i)
	}
	
	numbers := [] int {1,2,3,4,5,6}

	// iterating over collection
	for index , value := range numbers {
		fmt.Printf("Index: %d , Value: %d\n" , index , value) // %d is specific for numbers
	}

	for i:= 1 ; i <= 10 ; i++ {
		if(i % 2 == 0) {
			continue; // when compiler sees continue it does not move further rather than it moves back to top and continues from start.
		}
		fmt.Printf("Odd Number: %d\n" , i)

		if(i == 5) {
			break // break statement once it is run by the compiler , it exits from the loop and move forward...
		}
	}

	rows := 5
	
	// OuterLoop for nextLine
	for i := 1 ; i <= rows ; i++ {

		// inner loop to create spaces
		for j := 1 ; j <= rows - i ; j++ {
			fmt.Print(" ")
		}
		// inner loop to print astericks
		for k := 1 ; k <= 2 * i - 1 ; k++ {
			fmt.Print("*")
		}
		// to move to next line
		fmt.Println()
	}

	for i := range 10 {
		fmt.Println(10  - i)
	}
	fmt.Println("We have a LiftOff")

}