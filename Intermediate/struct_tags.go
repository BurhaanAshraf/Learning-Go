package main

import (
	"encoding/json"
	"fmt"
	"log"
)

/*
====================================================
TOPIC: STRUCT TAGS & JSON MARSHALING
====================================================

Struct Tags
-----------
Struct tags provide metadata about struct fields.

Example:
	FirstName string `json:"first_name"`

The json package reads these tags while
marshaling and unmarshaling.

Multiple tags can be used on the same field:

	FirstName string `json:"first_name" db:"firstn" xml:"first"`

Meaning:
	json package -> first_name
	database layer -> firstn
	xml package -> first

----------------------------------------------------

omitempty
----------
Removes field from JSON output when it contains
its zero value.

Examples:

	string -> ""
	int    -> 0
	bool   -> false
	slice  -> nil
	map    -> nil

Example:

	LastName string `json:"last_name,omitempty"`

If LastName == ""
JSON output will not contain last_name.

----------------------------------------------------

Ignoring Fields Completely
--------------------------

	json:"-"

Example:

	Password string `json:"-"`

Field will NEVER appear in JSON output,
even if it contains a value.

Difference:

omitempty
	Removes field only when it has zero value.

"-"
	Always removes field.

----------------------------------------------------

INTERVIEW QUESTIONS
-------------------

Q. What are struct tags?
A. Metadata attached to struct fields used by
   packages like json, xml, gorm, validator, etc.

Q. What does omitempty do?
A. Omits a field from JSON if it contains
   its zero value.

Q. Difference between omitempty and "-"?
A.
	omitempty -> conditional removal
	"-"       -> always ignored

====================================================
*/

type Persons struct {

	// Multiple struct tags.
	//
	// json:"first_name"
	// -> JSON key will be first_name
	//
	// db:"firstn"
	// -> Database column mapping
	//
	// xml:"first"
	// -> XML element mapping
	FirstName string `json:"first_name" db:"firstn" xml:"first"`

	// Omitted if LastName == ""
	LastName string `json:"last_name,omitempty"`

	// Included in JSON because no omitempty.
	Age int `json:"age"`

	// Example of ignored field:
	// Password string `json:"-"`
}

func structTags() {

	/*
	====================================================
	MARSHALING
	====================================================

	Go Struct
		↓
	json.Marshal()
		↓
	JSON ([]byte)

	Returns:
		[]byte
		error

	====================================================
	*/

	person := Persons{
		FirstName: "Jane",

		// LastName intentionally omitted.
		// Because of omitempty,
		// this field won't appear in JSON.
		// LastName: "Doe",

		Age: 50,
	}

	jsonData, err := json.Marshal(person)

	if err != nil {
		log.Fatal("Cannot Marshal Struct", err)
	}

	// Marshal returns []byte.
	// Convert to string for readable output.
	fmt.Println(string(jsonData))

	/*
	Expected Output:

	{
		"first_name":"Jane",
		"age":50
	}

	Notice:
		last_name is missing because
		LastName == "" and omitempty is used.
	*/

	/*
	====================================================
	1-MINUTE REVISION
	====================================================

	Struct Tag
		`json:"name"`

	Multiple Tags
		`json:"name" db:"col" xml:"tag"`

	omitempty
		Remove field if zero value.

	json:"-"
		Always ignore field.

	Marshal
		Struct -> JSON

	Unmarshal
		JSON -> Struct

	Marshal Returns
		[]byte, error

	====================================================
	*/
}