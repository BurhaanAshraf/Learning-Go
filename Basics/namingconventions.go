package main

import "fmt"


type EmployeeGoogle struct {
	FirstName string 
	LastName string
	Age int
}
type EmployeeApple struct {
	FirstName string 
	LastName string
	Age int
}


func namingConventions() {

	// PascalCase
	// Structs , enums , interfaces
	// CalculateArea , UserInfo

	//snake_case
	// Variables , constants and filenames
	//used_id , first_name

	//UPPERCASE
	// used to name Constants inside GO
	// This convention make sure constants stand out and their immutability is emphasized

	//camelCase
	// widely used to name variables 
	// example isValid , employeeID etc

	
	const MAXRETRIES = 5
	
	employeeID := 1001
	fmt.Println(employeeID)

}
