package main

import "fmt"

// ============================================================
// 🔹 1. Base Struct
// ============================================================

type person struct {
	name string
	age  int
}

// Method of person
func (p person) introduce() {
	fmt.Printf(
		"Hey my name is %s and I am %d years old.",
		p.name,
		p.age,
	)
	fmt.Println()
}


// ============================================================
// 🔹 2. Embedded Struct
// ============================================================

type employee struct {

	// Anonymous embedded struct
	// employee gets direct access to fields + methods of person

	person

	// If written like this:
	// personalInfo person
	// then it becomes a normal named field
	// access would be: employee.personalInfo.name

	employeeID string
	salary     float64
}


// ============================================================
// 🔹 3. Method Overriding-like Behavior
// ============================================================

// employee has its own introduce() method

func (e employee) introduce() {
	fmt.Printf(
		"Hey my name is %s and I am %d years old bearing ID %s and my salary is %.2f",
		e.name,
		e.age,
		e.employeeID,
		e.salary,
	)
	fmt.Println()
}


// ============================================================
// 🔹 4. Main Function
// ============================================================

func embedding() {

	employee01 := employee{
		person: person{
			name: "Burhaan",
			age:  20,
		},
		employeeID: "25bcs10462",
		salary:     250000,
	}

	// Direct access because person is anonymously embedded
	fmt.Println(employee01.name)

	fmt.Println(employee01.salary, employee01.employeeID)

	employee01.introduce()

	// Since employee embeds person:
	// - fields are promoted
	// - methods are promoted

	// But if employee defines its own method
	// with same name → employee's method is used first
}


// ============================================================
// 🔹 5. Important Concepts
// ============================================================

// Struct Embedding:
// One struct inside another struct

// Anonymous Embedding:
// person
// → direct field access: employee.name

// Named Field Embedding:
// personalInfo person
// → access using: employee.personalInfo.name

// Method Promotion:
// Embedded struct methods are available
// on outer struct automatically

// Method Overriding-like Behavior:
// If outer struct defines same method,
// outer method takes priority

// Go does NOT have true inheritance
// It uses composition through embedding