package main

import (
	"fmt"

	"github.com/Umbrasyl/pokedex/utils/commands"
)

func main() {

	fmt.Println("Welcome to the Pokedex REPL!")
	fmt.Println("Type 'help' for instructions or 'exit' to quit.")
	for {
		fmt.Printf(">> ")
		var input string
		fmt.Scanln(&input)
		// If the command is not found, print an error message
		if _, ok := commands.Commands[input]; !ok {
			fmt.Println("Unknown command - type 'help' for instructions")
			continue
		}
		if input == "exit" {
			commands.Commands[input].Callback()
			break
		}

		commands.Commands[input].Callback()
	}

}
