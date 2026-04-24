package main

import (
	"fmt"
)

// ============================================================
// 🔹 1. Struct Definitions
// ============================================================

type Person struct {
	firstName string
	lastName  string
	age       int

	address Address      // embedded struct (named field)
	PhoneHomeCell        // anonymous field (embedded type)
}

type Address struct {
	country  string
	city     string
	pincode  int
}

// Anonymous field struct
// Fields of this struct are promoted to Person
type PhoneHomeCell struct {
	home string
	cell string
}


// ============================================================
// 🔹 2. Main Function
// ============================================================

func Struct() {

	customer01 := Person{
		firstName: "Burhaan",
		lastName:  "Ashraf",
		age:       20,
		address: Address{
			country: "Korea",
			city:    "Seoul",
		},
		PhoneHomeCell: PhoneHomeCell{
			cell: "123456789",
			home: "987654321",
		},
	}

	// Assigning value later
	customer01.address.pincode = 123456

	fmt.Println(customer01)

	// Accessing fields
	fmt.Println(customer01.firstName)

	// Accessing nested struct field
	fmt.Println(customer01.address.pincode)

	// Accessing promoted fields from anonymous struct
	fmt.Println(customer01.cell)


	// ============================================================
	// 🔹 3. Struct Comparison
	// ============================================================

	compare01 := Person{
		firstName: "John",
		age:       21,
	}
	compare02 := Person{
		firstName: "John",
		age:       21,
	}

	fmt.Println(compare01 == compare02)
	// Structs are comparable if all fields are comparable


	// ============================================================
	// 🔹 4. Anonymous Struct
	// ============================================================

	user := struct {
		username string
		email    string
	}{
		username: "Tenz",
		email:    "Tenz@example.org",
	}

	fmt.Println(user.username, user.email)


	// ============================================================
	// 🔹 5. Methods on Struct
	// ============================================================

	fmt.Println(customer01.fullName())

	fmt.Printf("Before Age %d\n", customer01.age)

	customer01.incrementAgeByOne()

	fmt.Printf("After %d\n", customer01.age)
}


// ============================================================
// 🔹 6. Methods
// ============================================================

// Value Receiver (read-only behavior)
func (p Person) fullName() string {
	return p.firstName + " " + p.lastName
}

// Pointer Receiver (modifies original struct)
func (p *Person) incrementAgeByOne() {
	p.age += 1
}


// ============================================================
// 🔹 7. Important Concepts
// ============================================================

// Struct → collection of fields (like object)

// Embedded struct:
// Person has Address as a field → accessed via person.address.field

// Anonymous field (embedding):
// Fields are promoted → accessed directly (person.cell)

// Struct comparison:
// Allowed only if all fields are comparable

// Methods:
// Defined outside struct (Go design)
// Value receiver → does NOT modify original
// Pointer receiver → modifies original


// ============================================================
// 🔹 8. Go Design Rules
// ============================================================

// Structs and methods must be defined at package level
// Cannot define them inside functions

// Methods are declared separately from structs
// (unlike classes in OOP languages)

// Go separates:
// data (struct) and behavior (methods)