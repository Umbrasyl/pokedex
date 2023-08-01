package commands

import (
	"fmt"
	"time"

	"github.com/Umbrasyl/pokedex/utils/pokecache"
)

type cliCommand struct {
	name        string
	description string
	Callback    func()
}

var Commands map[string]cliCommand
var cache = pokecache.NewCache(5 * time.Minute)

func init() {
	Commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Show this help message",
			Callback: func() {
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
			Callback: func() {
				fmt.Println("Bye!")
			},
		},
		"map": {
			name:        "map",
			description: "List the next 20 locations from the PokeAPI with each call",
			Callback: func() {
				Mapper(20)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "List the previous 20 locations from the PokeAPI with each call",
			Callback: func() {
				Mapper(-20)
			},
		},
	}
}
