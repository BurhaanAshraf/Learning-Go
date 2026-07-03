package main

import (
	"fmt"
	"reflect"
)

type Human struct {
	Name string
	Age  int
}

func workingWithStructs() {
	person1 := Human{
		Name: "John",
		Age:  21,
	}

	// reflect.ValueOf(person1) wraps a COPY of person1 — not the original.
	// Copies are never addressable, so CanSet() is always false here.
	v := reflect.ValueOf(person1)
	fmt.Println("CanAddr:", v.CanAddr())
	fmt.Println("CanSet :", v.CanSet())

	for i := range v.NumField() {
		fmt.Printf("Field %d: %v\n", i, v.Field(i))
	}

	// To modify the real struct: pass a pointer, then .Elem() to reach
	// the value the pointer points to. That value IS addressable.
	v1 := reflect.ValueOf(&person1).Elem()

	fmt.Println("\nAfter pointer + Elem()")
	fmt.Println("CanAddr:", v1.CanAddr())
	fmt.Println("CanSet :", v1.CanSet())

	nameField := v1.FieldByName("Name")
	ageField := v1.FieldByName("Age")

	// CanSet() needs BOTH: addressable value + exported field.
	// Reflection never bypasses Go's normal lowercase/unexported rule.
	if nameField.CanSet() {
		nameField.SetString("Burhaan")
	}
	if ageField.CanSet() {
		ageField.SetInt(30)
	}

	fmt.Println("\nModified Struct")
	for i := range v1.NumField() {
		fmt.Printf("Field %d: %v\n", i, v1.Field(i))
	}
}

type Greeter struct{}

func (g Greeter) Greet(fname, lname string) string {
	return "Hello! " + fname + " " + lname
}

// Age has Kind int but Type main.Age — good example of Kind vs Type differing.
type Age int

func basics() {
	var x = 42

	v := reflect.ValueOf(x)
	t := reflect.TypeOf(x)

	fmt.Println("Value:", v)
	fmt.Println("Type :", t)

	// Kind = underlying category (int, string, struct...)
	// Type = the exact named type — differs from Kind for custom types.
	fmt.Println("Kind:", t.Kind())
	fmt.Println("Is Int:", t.Kind() == reflect.Int)
	fmt.Println("Is String:", t.Kind() == reflect.String)
	fmt.Println("Is Zero:", v.IsZero())

	fmt.Println()
	var custom Age = 25
	fmt.Println("Custom Type:", reflect.TypeOf(custom))        // main.Age
	fmt.Println("Custom Kind:", reflect.TypeOf(custom).Kind()) // int

	y := 10
	v1 := reflect.ValueOf(&y).Elem() // addressable -> SetInt works
	v2 := reflect.ValueOf(&y)        // this Value IS the *int, not the int itself

	fmt.Println()
	fmt.Println("V2 Type:", v2.Type())
	fmt.Println("Original Value:", v1.Int())
	v1.SetInt(18)
	fmt.Println("Modified Value:", v1.Int())

	// reflect.ValueOf on an interface unwraps it automatically — you get
	// the concrete underlying type straight away, never an "interface" Kind.
	fmt.Println()
	var itf interface{} = "Hello"
	v3 := reflect.ValueOf(itf)
	fmt.Println("V3 Type:", v3.Type())
	if v3.Kind() == reflect.String {
		fmt.Println("String Value:", v3.String())
	}
}

func workingWithMethods() {
	g := Greeter{}
	t := reflect.TypeOf(g)
	v := reflect.ValueOf(g)

	fmt.Println("Type:", t)

	var method reflect.Method
	for i := range t.NumMethod() {
		method = t.Method(i)
		fmt.Printf("Method %d: %s\n", i, method.Name)
	}

	// MethodByName on a Value (not a Type) returns a method already bound
	// to that receiver — no need to pass the receiver again when calling it.
	m := v.MethodByName(method.Name)

	// Call() takes []reflect.Value, not raw Go values, because each argument
	// needs to carry its runtime type info along with the value itself.
	results := m.Call([]reflect.Value{
		reflect.ValueOf("Alice"),
		reflect.ValueOf("Doe"),
	})
	fmt.Println("Greet Result:", results[0].String())
}

func Reflect() {
	basics()
	fmt.Println()
	workingWithStructs()
	fmt.Println()
	workingWithMethods()
}

/*
==================== NOTES ====================

Value vs Type
- reflect.Value -> the actual runtime data
- reflect.Type  -> metadata about it (fields, methods, name)

Addressability (the part that trips people up most)
- reflect.ValueOf(x)         -> a copy, CanSet() = false, always
- reflect.ValueOf(&x).Elem() -> addressable, CanSet() can be true
- Rule of thumb: if you didn't pass a pointer and call .Elem(), you
  cannot modify the value through reflection — full stop.

CanSet() requires BOTH:
1. an addressable value (see above)
2. an exported (capitalized) field
   Unexported fields stay unsettable even on an addressable struct —
   reflection respects normal Go visibility, it doesn't bypass it.

Kind vs Type
- Kind = underlying category (int, struct, slice, string...)
- Type = the exact named type (e.g. main.Age, whose Kind is still int)
- Use Kind() for generic switch/compare logic; Type() when you need
  the precise type identity.

Methods
- Type.Method(i)        -> just metadata (name, signature)
- Value.MethodByName(n) -> a real callable, already bound to its receiver
- Call() takes []reflect.Value so each argument carries its runtime type

When to actually reach for reflection
- Generic serializers (encoding/json), ORMs, validators, DI containers —
  situations where the type isn't known at compile time.
- Trade-offs: slower than direct code, no compile-time type safety, and
  it panics at runtime if you skip CanSet()/Kind() checks before acting.
  If a plain interface or a Go generic function (1.18+) solves the
  problem, prefer that over reflection.
==================================================
*/
