package main

import (
	"fmt"
	"math"
)

func arithmeticoperations() {

	const DIV int = 22 / 7 // this will give the integer value
	fmt.Println(DIV)

	const anotherDiv float64 = 22 / 7 // this will still give interger value besides we are using float datatype
	fmt.Println(anotherDiv) 

	const div float64 = 22 / 7.0 // this will give us float value because one the value should be float to get flaot as output
	fmt.Println(div)

	//OverFlow with signed integers

	var newInt int64 = 9223372036854775807 // this is the max value for int64 (signed)
	fmt.Println(newInt) 
	newInt = newInt + 1 // this will cause overflow
	fmt.Println(newInt)

	//Overflow with unsigned Integer

	var uMaxInt uint64 = 18446744073709551615 // this is the max value for unint64 (unsigned)
	fmt.Println(uMaxInt) 
	uMaxInt = uMaxInt + 1
	fmt.Println(uMaxInt)


	//Underflow with floating point numbers

	var smallFloat float64 = 1.0e-323

	fmt.Println(smallFloat)

	smallFloat = smallFloat / math.MaxFloat64 // this is underflow , because due to loss of precision it rounded up to 0

	fmt.Println(smallFloat)
}
