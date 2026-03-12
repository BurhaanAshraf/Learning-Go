package main

import "fmt"

var globalScope = "This is global scope variable" // we have declared a package level variables and we cannot create it using goofer

func main() {
	var age int

	age = 10

	fmt.Println(age)


	var name string // here we have not initialized name with data that is why we have to provide it with datatype
	name = "Burhaan"

	fmt.Println(name)
	var her = "Twaha"
	fmt.Println(her) // here we can omit writing the type because we have initialised it with string data


	count := 10
	fmt.Println(count)
	hisName := "Burhaan" //  here we are not using the var keyboard and also we are not declaring the data...   :=
	fmt.Println(hisName)

	// we cannot use var count int := 10 , it will throw an error

	
	// default values

	// Numeric Types --> 0

	// Boolean Types --> False

	// String --> ""

	// pointers , slices , maps , structs and functions --> nil

	{
			// block scoping of variables inside Go

			newBlock := "This is inside Block and I am trying to do block scoping"

			fmt.Println(newBlock)
	}

	// fmt.Println(newBlock) // this will thrown an error because --> we are trying to access newBlock outside its block scope 
	

}



