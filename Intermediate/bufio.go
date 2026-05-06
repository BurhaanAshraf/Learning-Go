package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


func BufIo() {

	reader := bufio.NewReader(strings.NewReader("Hello, BufIo PackageEE!\n")) 
	
	// bufIo.NewReader accepts NewReader as a parameter and then buffers the exiting Reader
	// it returns the bufIO.Reader instance

	// Reading byte slice

	data := make([]byte ,20)

	n , err :=reader.Read(data) // data is the destination we will read from reader to data 

	if err != nil {
		panic("Cannot Read")
	}
	fmt.Printf("Read: %d bytes: %s\n" , n , data[:n])

	newStr , err := reader.ReadString('\n') // reading starts from the point where it stops in that case it is afte 20 bytes --> data

	if err != nil {
		panic("Cannot Read")
	}
	fmt.Println("Read String:", newStr)


	// Bufio.Writer --> here we declare the target first and then source

	writer := bufio.NewWriter(os.Stdout)

	// Writing byte slice

	data = ([]byte("Hello , bufio package!"))

	m , err := writer.Write(data)

	if err != nil {
		panic("Error Writing")
	}
	fmt.Println(m ,"bytes")
	//all the data that is written to writer is stored in internal buffer and it is not immediately return to writer that we pass as argument

	// Flush the buffer to ensure all data is written to target writer

	err = writer.Flush()
	fmt.Println()

	if err != nil {
		panic("Error Flushing Writer")
	}
	// everytime we need to pass data to writer we need to manually flush data to writer

	// Write String

	newWriter := bufio.NewWriter(os.Stdout)
	writeString := "This is a string\n"

	a , err := newWriter.WriteString(writeString)

	if err != nil {
		panic("Cannot Write String")
	}
	fmt.Println(a) // this will print number of bytes that were wrote
	err = newWriter.Flush() 
	if err != nil {
		panic("Cannot Flush")
	}
	



	

	










}