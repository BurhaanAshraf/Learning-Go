package main

import "fmt"


func main() {

	num := 42
	fmt.Printf("%05d\n" , num) // this five will only be active when we have less than 5 digits

	message := "Hello"
	fmt.Printf("|%10s|\n" , message) // here we are fixing the width to minimum 10 , if  < 10 then we will have leading spaces
	fmt.Printf("|%-10s|\n" , message ) // here we are fixing the width to minimum 10 , if < 10 then we will have trailing spaces

	// if we use backtick then it would be a RAW string which means it takes everything into string literal!

	message2 := "Hello \nWorld" // this will use escape sequence as new line
	message3 := `Hello \nWorld` // this will use escape sequence as string literal

	fmt.Println(message2)
	fmt.Println(message3)




}