package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Umbrasyl/pokedex/utils/commands"
)

func main() {

	fmt.Println("Welcome to the Pokedex REPL!")
	fmt.Println("Type 'help' for instructions or 'exit' to quit.")
	for {
		fmt.Printf(">>> ")
		// Read the input from the user
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		// Remove the newline character from the input
		input = input[:len(input)-1]

		inputSlice := strings.Split(input, " ")
		command := inputSlice[0]
		args := inputSlice[1:]

		if cmd, ok := commands.Commands[command]; ok {
			if command == "exit" {
				cmd.Callback(args)
				break
			}
			cmd.Callback(args)
		} else {
			fmt.Println("Unknown command:", command)
		}
	}

}
