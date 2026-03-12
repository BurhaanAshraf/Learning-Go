package main

import ( // import () when we have have to import multiple packages no need to seperate by commas.
	"fmt"
	foo "net/http" // named import
)


func server() {
	fmt.Println("Hello GO Wolrd")
	 resp , err := foo.Get("https://jsonplaceholder.typicode.com/posts/1")
	 if(err != nil) {
		fmt.Println("Error" ,err)
	 }
	 defer resp.Body.Close()
	 fmt.Println("Response is " , resp.Status)
	 }