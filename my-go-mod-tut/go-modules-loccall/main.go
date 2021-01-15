package main

import (
	"fmt"

	rv "github.com/samueldsr/sam-go-mod/morestrings"
)

func main() {
	msg := "hello golang, import local packages/subpackages"
	fmt.Println("Original message: ", msg)

	msg = rv.ReverseRunes(msg)
	fmt.Println("Reversed message: ", msg)

}
