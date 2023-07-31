package main

import "fmt"

func main() {
	// Buidl a REPL
	fmt.Println(">> Welcome to the Pokedex REPL!")
	fmt.Println(">> Type 'help' for instructions or 'exit' to quit.")
	for {
		fmt.Printf(">> ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "exit":
			fmt.Println("Bye!")
			return
		case "help":
			fmt.Println("You can use the program by typing commands when prompted.")
			fmt.Println("The commands are:")
			fmt.Println("  help   - show this help message")
			fmt.Println("  exit   - exit the program")
		default:
			fmt.Println("Unknown command")
		}
	}
}
