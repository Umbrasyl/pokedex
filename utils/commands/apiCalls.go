package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var Mapper func(int) = func() func(int) {
	// The function use closures to keep track of the current page
	currPagination := -20
	return func(pagination int) {
		currPagination += pagination
		if currPagination < 0 {
			currPagination = 0
		}
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d&limit=20", currPagination)
		// Try to get the data from the cache. Get method returns two values: the data and a boolean.
		// The boolean is true if the data is found in the cache, false otherwise.
		data, ok := Cache.Get(url)
		if ok {
			// Data is a []byte, so we need to unmarshal it into the locations struct
			Cache.UpdateTime(url)
		} else {
			// If the data is not found in the cache, fetch it from the API
			fetchedData, err := getResource(url)
			data = fetchedData
			if err != nil {
				fmt.Println(err)
				return
			}
			// Add the fetchedData to the cache
			Cache.Add(url, fetchedData)
		}
		var locations Locations
		err := json.Unmarshal(data, &locations)
		if err != nil {
			fmt.Println("Couldn't get the location: ", err)
			return
		}
		printLocations(locations.Results)
	}
}()

func printLocations(locations []struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}) {
	for _, location := range locations {
		// Split the location name by dashes and capitalize each word
		// e.g. "seafoam-islands" -> "Seafoam Islands"
		var locationName string
		for _, word := range strings.Split(location.Name, "-") {
			locationName += cases.Title(language.English).String(word) + " "
		}
		fmt.Println(locationName)
	}
}

// getLocations should fetch the data and return it as a []byte
func getResource(url string) ([]byte, error) {
	// Make a GET request to the API
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// Close the response body when the function returns
	defer resp.Body.Close()
	// Read the response body into a byte slice
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Explore should print the all the pokemon that can be found in the given area
func Explore(areaName string) {
	// The url looks like https://pokeapi.co/api/v2/location-area/{id or name}/
	// First check the cache for the data
	fmt.Println("Exploring", areaName)
	url := "https://pokeapi.co/api/v2/location-area/" + areaName + "/"
	data, ok := Cache.Get(url)
	if ok {
		// If the data is found in the cache, unmarshal it into the area struct
		Cache.UpdateTime(url)
	} else {
		// If the data is not found in the cache, fetch it from the API
		fetchedData, err := getResource(url)
		data = fetchedData
		if err != nil {
			fmt.Println(err)
			return
		}
		// Add the fetchedData to the cache
		Cache.Add(url, fetchedData)
	}
	var area Area
	err := json.Unmarshal(data, &area)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, encounter := range area.PokemonEncounters {
		fmt.Println(" - " + encounter.Pokemon.Name)
	}
}

func Catch(pokemonName string) {
	// The url looks like https://pokeapi.co/api/v2/pokemon/{id or name}/
	// First check the cache for the data
	fmt.Println("Throwing a ball at", pokemonName)
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName + "/"
	data, ok := Cache.Get(url)
	if ok {
		// If the data is found in the cache, unmarshal it into the pokemon struct
		Cache.UpdateTime(url)
	} else {
		// If the data is not found in the cache, fetch it from the API
		fetchedData, err := getResource(url)
		data = fetchedData
		if err != nil {
			fmt.Println(err)
			return
		}
		// Add the fetchedData to the cache
		Cache.Add(url, fetchedData)
	}
	var pokemon Pokemon
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Catching will be based on the pokemon's baseExperience and random chance
	chanceToCatch := ((rand.Float64()*5 + 5) / 10) * (100 / float64(pokemon.BaseExperience))
	rollToCatch := rand.Float64()
	if rollToCatch <= chanceToCatch {
		fmt.Println(pokemonName, "was caught!")
		caughtPokemon[pokemonName] = pokemon
	} else {
		fmt.Println(pokemonName, "got away!")
	}
}

func Inspect(pokemonName string) {
	// This function will look at the caughtPokemon map and print the pokemon's stats if it exists
	pokemon, ok := caughtPokemon[pokemonName]
	if !ok {
		fmt.Println("You don't have", pokemonName)
		return
	}
	fmt.Println("Name: ", pokemon.Name)
	fmt.Println("Height: ", pokemon.Height)
	fmt.Println("Weight: ", pokemon.Weight)
	fmt.Println("Stats: ")
	fmt.Println(" - HP: ", pokemon.Stats[0].BaseStat)
	fmt.Println(" - Attack: ", pokemon.Stats[1].BaseStat)
	fmt.Println(" - Defense: ", pokemon.Stats[2].BaseStat)
	fmt.Println(" - Special Attack: ", pokemon.Stats[3].BaseStat)
	fmt.Println(" - Special Defense: ", pokemon.Stats[4].BaseStat)
	fmt.Println(" - Speed: ", pokemon.Stats[5].BaseStat)
	fmt.Println("Types: ")
	for _, pokemonType := range pokemon.Types {
		fmt.Println(" - ", pokemonType.Type.Name)
	}
}

func Pokedex() {
	// This function will print the names of all the pokemon that have been caught
	fmt.Println("Pokedex:")
	for pokemonName := range caughtPokemon {
		fmt.Println(" - ", pokemonName)
	}
}
