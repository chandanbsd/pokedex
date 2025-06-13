package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	res := strings.Fields(text)

	values := []string{}

	for _, s := range res {
		values = append(values, strings.ToLower(s))
	}

	fmt.Print(len(values))
	return values
}
