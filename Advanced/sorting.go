package main

import (
	"fmt"
	"sort"
)

// PATTERN: "function as a sort criterion" — lets you sort by any rule
// (age, name, custom logic) without writing a new type for each one.
// This is the classic Go idiom from the standard library's sort.Interface,
// wrapped so the comparison logic itself becomes pluggable.

type Person struct {
	name string
	age  int
}

// By is a comparison function: returns true if p1 should come before p2.
// Defining it as a named func type lets us attach a Sort method to it directly.
type By func(p1, p2 *Person) bool

// PersonSorter adapts a slice + a By function into something, sort.Sort() can use.
// sort.Sort() needs three methods: Len, Less, Swap — this struct provides them
// by delegating the actual comparison to the By function.
type PersonSorter struct {
	people []Person
	by     func(p1, p2 *Person) bool
}

func (s *PersonSorter) Len() int {
	return len(s.people)
}

func (s *PersonSorter) Less(i, j int) bool {
	return s.by(&s.people[i], &s.people[j]) // delegate comparison to whichever By func was passed in
}

func (s *PersonSorter) Swap(i, j int) {
	s.people[i], s.people[j] = s.people[j], s.people[i]
}

// Sort is a method on By itself — this is what makes By(ageAsc).Sort(people) work.
// It wraps the slice + comparison func into a PersonSorter and hands it to sort.Sort,
// which calls Len/Less/Swap internally until the slice is ordered.

// By is NOT doing the sorting.

// By is NOT storing data.

// By is NOT changing the comparator.

// By is just a named function type so we can attach
// a Sort() method to a comparison function.

func (by By) Sort(people []Person) {
	sort.Sort(&PersonSorter{people: people, by: by})
}

func sorting() {
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Catherine", 15},
	}
	fmt.Println("Before sorting:", people)

	// Each of these is just a func(p1, p2 *Person) bool — the comparison rule.
	ageAsc := func(p1, p2 *Person) bool { return p1.age < p2.age }
	ageDesc := func(p1, p2 *Person) bool { return p1.age > p2.age }
	byName := func(p1, p2 *Person) bool { return p1.name < p2.name }

	By(ageAsc).Sort(people)
	// A function can be converted to a named function type if their signatures are identical.
	// That's why ageAsc can become a by, and once it becomes a by,
	// it gains access to the methods attached to by (like sort()).
	fmt.Println("Sorted by age asc:", people)

	By(byName).Sort(people)
	fmt.Println("Sorted by name:", people)

	By(ageDesc).Sort(people)
	fmt.Println("Sorted by age desc:", people)

	// sort.Slice — the simpler, modern alternative when you don't need to reuse
	// the comparator elsewhere. No interface, no extra types. Just a slice + a Less func.
	words := []string{"Banana", "Apple", "Grapes", "Cherry", "Guava", "Watermelon", "Pear", "Orange"}
	sort.Slice(words, func(i, j int) bool {
		return words[i][len(words[i])-1] < words[j][len(words[j])-1] // compare last character of each word
	})
	fmt.Println("Sorted by last character:", words)
}
