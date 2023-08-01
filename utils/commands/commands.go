package commands

import (
	"fmt"
	"time"

	"github.com/Umbrasyl/pokedex/utils/pokecache"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(args []string)
}

var Commands map[string]cliCommand
var Cache = pokecache.NewCache(2 * time.Minute)

func init() {
	Commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Show this help message",
			Callback: func(args []string) {
				fmt.Println("You can use the program by typing commands when prompted.")
				fmt.Println("The commands are:")
				for _, command := range Commands {
					fmt.Printf("%s - %s\n", command.name, command.description)
				}
			},
		},
		"exit": {
			name:        "exit",
			description: "Exit the program",
			Callback: func(args []string) {
				fmt.Println("Bye!")
			},
		},
		"map": {
			name:        "map",
			description: "List the next 20 locations from the PokeAPI with each call",
			Callback: func(args []string) {
				Mapper(20)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "List the previous 20 locations from the PokeAPI with each call",
			Callback: func(args []string) {
				Mapper(-20)
			},
		},
		"explore": {
			name:        "explore",
			description: "Give a location name to print all the pokemon in that area",
			Callback: func(args []string) {
				if len(args) == 0 {
					fmt.Println("Please provide a location name")
					return
				}
				Explore(args[0])
			},
		},
	}
}
