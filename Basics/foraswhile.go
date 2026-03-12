package main

import "fmt"


func forasawhile() {

	i := 1
	for i <= 5 {
		fmt.Println("Iteration:",i)
		i++
	}

	sum := 0
	for {
		fmt.Println("Sum:", sum)
		sum += 10
		if(sum == 50) {
			break
		}
	}

	num := 1
	for num <= 10 {
		if num % 2 == 0 {
			num++
			continue
		}
		fmt.Println("Odd Number", num)
		num++
	}

	



}
