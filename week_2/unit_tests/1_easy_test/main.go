package main

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
)

func main() {
	isBrain := gofakeit.Bool()
	if areYouStupid(isBrain) {
		fmt.Println("You are stupid")
		return
	}
	fmt.Println("You are not stupid")
}

func areYouStupid(isBrain bool) bool {
	if isBrain {
		return true
	}
	return false
}

func TestAreYouStupid() {
	fmt.Println("TestAreYouStupid true test")
	if !areYouStupid(true) {
		fmt.Println("Test failed")
		return
	}

	fmt.Println("Test passed")

	fmt.Println("TestAreYouStupid false test")
	if areYouStupid(false) {
		fmt.Println("Test failed")
		return
	}

	fmt.Println("Test passed")
}
