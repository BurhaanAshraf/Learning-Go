package main

import (
	"encoding/json" // Package used for JSON encoding and decoding.
	"fmt"           // Used for printing output.
	"log"           // Used for logging errors.
)

/*
=====================================================
REVISION NOTES
=====================================================

JSON Marshalling
----------------
Go Struct  ---> JSON

json.Marshal()

JSON UnMarshalling
------------------
JSON ---> Go Struct / Map

json.Unmarshal()

Interview Question:
Q: Why are struct fields capitalized?
A: Only exported (capitalized) fields are accessible
   by the json package.
=====================================================
*/

type Human struct {

	// JSON package looks at these tags while
	// converting struct to JSON.

	Name string `json:"name"`

	// omitempty removes field if value is zero value.
	// int -> 0
	// string -> ""
	// bool -> false
	Age int `json:"age,omitempty"`

	Email string `json:"Email"`

	// Nested Struct
	Address ADDRESS `json:"address"`
}

/*
=====================================================
NESTED STRUCT

Human
 └── Address
      ├── City
      ├── State
      └── Pincode

Will become nested JSON object.
=====================================================
*/
type ADDRESS struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Pincode int    `json:"pincode"`
}

/*
Employee struct used for UnMarshalling JSON.

JSON keys must match JSON tags for automatic mapping.
*/
type Employee struct {
	FullName string  `json:"full_name"`
	EmpID    string  `json:"emp_ID"`
	Age      int     `json:"age"`
	Address  ADDRESS `json:"Address"`
}

func JSON() {

	/*
		=====================================================
		MARSHALLING EXAMPLE 1
		=====================================================

		Converting Struct -> JSON
	*/

	person := Human{
		Name:  "John",
		Email: "john@example.com",
	}

	jsonData, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}

	// Convert []byte into string before printing.
	fmt.Println(string(jsonData))

	/*
		Output:

		{
			"name":"John",
			"Email":"john@example.com"
		}

		Notice:
		Age is missing because of omitempty.
	*/

	person1 := Human{
		Name:  "Jane",
		Age:   21,
		Email: "jane@example.com",

		Address: ADDRESS{
			City:    "Bangalore",
			State:   "KA",
			Pincode: 123456,
		},
	}

	person1Data, err := json.Marshal(person1)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(person1Data))

	/*
		=====================================================
		UNMARSHALLING
		=====================================================

		JSON -> Struct

		json.Unmarshal(
			jsonBytes,
			&structVariable
		)
	*/

	jsonObjData := `{
		"full_name"  : "Jenny Doe",
		"emp_ID" : "0009",
		"age" : 30,
		"address" : {
			"city" : "Bandra",
			"state" : "Mumbai",
			"pincode" : 567890
		}
	}`

	var employeeFromJson Employee

	/*
		Struct currently has zero values.

		String -> ""
		Int    -> 0
		Bool   -> false
	*/
	fmt.Println(employeeFromJson.FullName)

	// Decode JSON into struct.
	err = json.Unmarshal([]byte(jsonObjData), &employeeFromJson)

	if err != nil {
		panic(err)
	}

	/*
		Printing entire struct.

		Note:
		Field names are not printed,
		only values.
	*/
	fmt.Println(employeeFromJson)

	fmt.Println(employeeFromJson.Age)
	fmt.Println(employeeFromJson.FullName)

	/*
		=====================================================
		HANDLING LIST / ARRAY OF OBJECTS
		=====================================================

		[]ADDRESS means Slice of ADDRESS structs.
	*/

	ListOfAddress := []ADDRESS{
		{City: "New York", State: "NY", Pincode: 123456},
		{City: "Washington DC", State: "WDC", Pincode: 568373},
		{City: "Bangalore", State: "KA", Pincode: 781298},
		{City: "Bandra", State: "Mumbai", Pincode: 751234},
		{City: "Chandigarh", State: "Haryana", Pincode: 981456},
		{City: "Patna", State: "Bihar", Pincode: 347653},
	}

	fmt.Println(ListOfAddress)

	// Convert slice into JSON array.
	newData, err := json.Marshal(ListOfAddress)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(newData))

	/*
		=====================================================
		UNKNOWN JSON STRUCTURE
		=====================================================

		When structure is not known beforehand,
		use:

		map[string]any

		Key   -> string
		Value -> any

		Equivalent to:

		JSON Object
		   ↓
		map[string]any
	*/

	unknownData := `{
		"full_name" : "Cole",
		"age" : 23,
		"address" : {
			"city" : "Bangkok",
			"country" : "Thailand"
		}
	}`

	// Create empty map.
	data := make(map[string]any)

	err = json.Unmarshal([]byte(unknownData), &data)

	if err != nil {
		log.Fatalln("Cannot UnMarshal JSON", err)
	}

	fmt.Println("Decoded JSON", data)

	/*
		Accessing top-level keys.

		data =
		{
			"full_name": "Cole",
			"age": 23
		}
	*/
	fmt.Println("Decoded Full Name", data["full_name"])
	fmt.Println("Decoded Age", data["age"])

	/*
		=====================================================
		TYPE ASSERTION
		=====================================================

		data["address"]

		returns type: any

		But we know internally it is:

		map[string]any

		So we tell Go:

		"Trust me, this value is a map."

		Syntax:

		value.(ActualType)

		This is called Type Assertion.
	*/

	fmt.Println(
		"Decoded City From Address",
		data["address"].(map[string]any)["city"],
	)

	/*
		=====================================================
		MEMORY TRICK
		=====================================================

		Marshal   = Go -> JSON

		Unmarshal = JSON -> Go

		M = Make JSON
		U = Use JSON

		Marshal   -> Struct to JSON
		Unmarshal -> JSON to Struct
		=====================================================
	*/
}