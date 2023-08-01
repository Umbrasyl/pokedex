package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type area struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

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
			fetchedData, err := getLocations(url)
			data = fetchedData
			if err != nil {
				fmt.Println(err)
				return
			}
			// Add the fetchedData to the cache
			Cache.Add(url, fetchedData)
		}
		var locations locations
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
func getLocations(url string) ([]byte, error) {
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
		fetchedData, err := getArea(url)
		data = fetchedData
		if err != nil {
			fmt.Println(err)
			return
		}
		// Add the fetchedData to the cache
		Cache.Add(url, fetchedData)
	}
	var area area
	err := json.Unmarshal(data, &area)
	if err != nil {
		fmt.Println(err)
		return
	}
	printArea(area)
}

func printArea(area area) {
	for _, encounter := range area.PokemonEncounters {
		fmt.Println(" - " + encounter.Pokemon.Name)
	}
}

func getArea(url string) ([]byte, error) {
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
