package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func main() {

	var word string
	scanner := bufio.NewScanner(os.Stdin)

	var cleanedSlice []string
	for true {
		fmt.Printf("Pokedex > ")
		isComplete := scanner.Scan()

		if isComplete {
			word = scanner.Text()
			cleanedSlice = cleanInput(word)

			if len(cleanedSlice) > 0 {
				commandStr := cleanedSlice[0]

				command, ok := commands[commandStr]

				if !ok {
					fmt.Println("Unknown command")
				} else {
					command.callback()
				}
			}

		}

	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {

	fmt.Println("Welcome to the Pokedex!\nUsage: \n")

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func cleanInput(text string) []string {
	res := strings.Fields(text)

	values := []string{}

	for _, s := range res {
		values = append(values, strings.ToLower(s))
	}

	return values
}
