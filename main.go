package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	var word string

	var cleanedSlice []string

	for true {
		fmt.Printf("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)

		isComplete := scanner.Scan()

		if isComplete {
			word = scanner.Text()
			cleanedSlice = cleanInput(word)

			if len(cleanedSlice) > 0 {
				fmt.Printf("Your command was: %s\n", cleanedSlice[0])
			}

		}

	}
}

func cleanInput(text string) []string {
	res := strings.Fields(text)

	values := []string{}

	for _, s := range res {
		values = append(values, strings.ToLower(s))
	}

	return values
}
