package main

import (
	"fmt"
	"math/rand"
	"time"
)

func guessinggame() {

	source := rand.NewSource(time.Now().UnixNano()) 
		random := rand.New(source)

		target := random.Intn(100) + 1

		fmt.Println("Welcome to the guessing game!")

		var guess int 

		tries := 10
		
		// here we are using for loop as a while loop
		for tries > 0 {
			fmt.Println("Enter your number")
			fmt.Scanln(&guess)

			if(guess == target) {
				tries--
				fmt.Println("Congratulations! You guessed it right")
				fmt.Printf("You guessed in %d tries\n" , 10 - tries)
				break
			}else if guess < target {
				tries--
				fmt.Println("Too low! Try guessing higher number")
				fmt.Println("Tries Left", tries)
			}else {
				tries--
				fmt.Println("Too High! Try guess low")
				fmt.Println("Tries Left", tries)
			}
			
		}

}
