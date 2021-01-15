package main

import (
	"fmt"

	"github.com/samueldsr/sam-go-mod/morestrings"
)

func main() {

	msg := "hello world, go"
	fmt.Println("Original message: ", msg)

	msg = morestrings.ReverseRunes(msg)
	fmt.Println("Reversed message: ", msg)
}
