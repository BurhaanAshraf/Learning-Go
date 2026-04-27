package main

import "fmt"

// ============================================================
// 🔹 Generics Definition
// ============================================================

// Generics allow us to write reusable code
// that works with multiple data types
// without rewriting the same logic.

// ============================================================
// 🔹 Generic Syntax
// ============================================================

// Generic Function Syntax:

// func functionName[T constraint](params) returnType { }

// Example:
// func swap[T any](a, b T) (T, T)

// T → type parameter
// any → T can be any type


func swap[T any](a, b T) (T, T) {
	return b, a
}


// ============================================================
// 🔹 Generic Struct
// ============================================================

// Generic Struct Syntax:

// type StructName[T constraint] struct { }

type Stack[T any] struct {
	elements []T
}


// ============================================================
// 🔹 Stack Methods
// ============================================================

// Push → adds element to stack

func (s *Stack[T]) push(element T) {
	s.elements = append(s.elements, element)
}


// Pop → removes last element (LIFO)

func (s *Stack[T]) pop() T {

	// If stack is empty,
	// return zero value of that type

	if len(s.elements) == 0 {
		var zero T
		return zero
	}

	element := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]

	return element
}


// isEmpty → checks if stack is empty

func (s *Stack[T]) isEmpty() bool {
	return len(s.elements) == 0
}


// printAll → prints all elements

func (s *Stack[T]) printAll() {
	for _, i := range s.elements {
		fmt.Print(i, " ")
	}
	fmt.Println()
}


// ============================================================
// 🔹 Main Function
// ============================================================

func Generics() {

	// Generic function with int

	a, b := swap(10, 20)
	fmt.Println(a, b)


	// Generic function with string

	c, d := swap("Jane", "John")
	fmt.Println(c, d)


	// ========================================================
	// Generic Stack with int
	// ========================================================

	myStack := Stack[int]{}

	myStack.push(1)
	myStack.push(2)
	myStack.push(3)

	myStack.printAll()

	fmt.Println(myStack.pop())

	myStack.printAll()

	fmt.Println("Is stack Empty:", myStack.isEmpty())


	// ========================================================
	// Generic Stack with string
	// ========================================================

	stringStruct := Stack[string]{}

	stringStruct.push("Burhaan")
	stringStruct.push("Bottle")

	stringStruct.printAll()

	fmt.Println(stringStruct.pop())

	stringStruct.printAll()
}


// ============================================================
// 🔹 Quick Revision Points
// ============================================================

// [T any]

// T → placeholder for actual type

// any → accepts all types

// Stack[int]
// Stack[string]

// var zero T
// gives zero value of that type

// int → 0
// string → ""

// Generics = write once, use for many types