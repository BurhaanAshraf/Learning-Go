package main

import (
	"encoding/xml"
	"fmt"
	"log"
)

type persons struct {

	// XMLName specifies the root XML tag.
	//
	// Without XMLName, Go uses the struct name as
	// the root element during marshaling.
	//
	// Useful when XML schema requires a specific
	// root tag name.
	XMLName xml.Name `xml:"person"`

	// Maps Name field to <name> XML element.
	Name string `xml:"name"`

	// omitempty removes this field from generated XML
	// if it contains the zero value for its type.
	//
	// int zero value = 0
	Age int `xml:"age,omitempty"`

	// Maps Email field to <email> element.
	Email string `xml:"email"`

	// Nested structs automatically become nested XML.
	//
	// This is the recommended way to represent
	// hierarchical XML data.
	Address Addresses `xml:"address"`
}

type Addresses struct {

	// Omitted if empty string because of omitempty.
	City string `xml:"city,omitempty"`

	Country string `xml:"country"`
}

func XML() {

	Person := persons{
		Name:  "John",
		Age:   30,
		Email: "john@example.com",

		Address: Addresses{
			City:    "New York",
			Country: "US",
		},
	}

	// ------------------------------------------------
	// MARSHALING
	// ------------------------------------------------
	//
	// Converts:
	//
	// Go Struct
	//     ↓
	// XML
	//
	// Returns:
	// []byte -> XML data
	// error  -> if conversion fails
	//
	// Common Use Cases:
	// - API responses
	// - XML file generation
	// - Data exchange with legacy systems
	//
	xmlData, err := xml.Marshal(Person)

	if err != nil {
		log.Fatal("Cannot Marshal Data", err)
	}

	// Suppressing unused variable warning.
	_ = xmlData

	// ------------------------------------------------
	// MARSHALINDENT
	// ------------------------------------------------
	//
	// Same as Marshal but generates human-readable XML.
	//
	// Preferred when:
	// - Debugging
	// - Logging
	// - Saving XML files for humans
	//
	// Parameters:
	// data   -> object to marshal
	// prefix -> added before each line
	// indent -> indentation characters
	//
	xmlData1, err := xml.MarshalIndent(Person, "", "  ")

	if err != nil {
		log.Fatal("Cannot Marshal Data", err)
	}

	fmt.Println(string(xmlData1))

	// ------------------------------------------------
	// UNMARSHALING
	// ------------------------------------------------
	//
	// Converts:
	//
	// XML
	//   ↓
	// Go Struct
	//
	// Important:
	// Must pass a pointer because Unmarshal
	// needs to write decoded values into
	// the target variable.
	//
	XMLData := `<person><name>John</name><age>25</age></person>`

	var personXML persons

	err = xml.Unmarshal([]byte(XMLData), &personXML)

	if err != nil {
		log.Fatal("Cannot UnMarshal XML", err)
	}

	// ------------------------------------------------
	// UNMARSHALING NESTED XML
	// ------------------------------------------------
	//
	// If XML structure matches struct tags,
	// encoding/xml automatically maps nested
	// elements into nested structs.
	//
	xmlRaw := `<person><name>Burhaan</name><age>20</age><email>burhaan@example.com</email><address><city>Bangalore</city><country>IN</country></address></person>`

	var personxml persons

	err = xml.Unmarshal([]byte(xmlRaw), &personxml)

	if err != nil {
		log.Fatal("Cannot UnMarshal XML", err)
	}

	fmt.Println(personxml)

	// ------------------------------------------------
	// XMLName FIELD
	// ------------------------------------------------
	//
	// After unmarshaling, XMLName stores
	// information about the root XML element.
	//
	// Useful when:
	// - Validating XML documents
	// - Working with multiple XML schemas
	//
	// Example:
	//
	// personxml.XMLName.Local
	//
	// Returns root tag name.
	//

	// fmt.Println(personxml.XMLName.Local)

	// ------------------------------------------------
	// XML ATTRIBUTES
	// ------------------------------------------------
	//
	// XML supports:
	//
	// 1. Elements
	// 2. Attributes
	//
	// Elements contain data as child nodes.
	//
	// Attributes store metadata inside tags.
	//
	// Go uses:
	//
	// `xml:"field,attr"`
	//
	// to convert a struct field into an attribute.
	//
	book := Book{
		ISBN:   "467-736-123-7654",
		Title:  "GOLANG",
		Author: "Burhaan",
	}

	bookData, err := xml.MarshalIndent(book, "", "  ")

	if err != nil {
		log.Fatal("Cannot Marshal XML Data", err)
	}

	fmt.Println(string(bookData))
}

type Book struct {

	// Sets root XML element name.
	//
	// NOTE:
	// Conventionally use XMLName,
	// not XMlName.
	XMLName xml.Name `xml:"book"`

	// attr converts field into an XML attribute.
	//
	// Useful for identifiers and metadata.
	//
	// Interview:
	// Difference between attribute and element?
	//
	// Attributes:
	// - Compact
	// - Metadata
	//
	// Elements:
	// - Main content
	// - Can contain nested structures
	//
	ISBN string `xml:"isbn,attr"`

	// Another attribute.
	Title string `xml:"title,attr"`

	// Normal XML element.
	Author string `xml:"author"`
}

/*
====================================================
REVISION NOTES
====================================================

Marshal
	Go Struct -> XML

MarshalIndent
	Go Struct -> Formatted XML

Unmarshal
	XML -> Go Struct

XMLName
	Controls root XML tag.

omitempty
	Skips field if it contains zero value.

attr
	Converts field into XML attribute.

Nested Structs
	Become nested XML automatically.

Unmarshal Requirement
	Always pass pointer.

Common Interview Questions

1. Why use XMLName?
2. Difference between Marshal and MarshalIndent?
3. Why does Unmarshal require a pointer?
4. Difference between XML elements and attributes?
5. What does omitempty do?

====================================================
*/