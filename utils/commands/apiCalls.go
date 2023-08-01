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

var Mapper func(int) = func() func(int) {
	// The function use closures to keep track of the current page
	currPagination := -20
	return func(pagination int) {
		currPagination += pagination
		if currPagination < 0 {
			currPagination = 0
		}
		locations, err := getLocations(currPagination)
		if err != nil {
			fmt.Println("Error fetching locations:", err)
			return
		}
		for _, location := range locations {
			// Split the location name by dashes and capitalize each word
			// e.g. "seafoam-islands" -> "Seafoam Islands"
			var locationName string
			for _, word := range strings.Split(location, "-") {
				locationName += cases.Title(language.English).String(word) + " "
				fmt.Println(locationName)
			}
		}
	}
}()

func getLocations(currPagination int) ([]string, error) {
	// Fetch the locations using
	// https://pokeapi.co/api/v2/location/?offset=currPagination&limit=20

	res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location/?offset=%d&limit=20", currPagination))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// Read the response body. It has the structure of the locations struct
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// Unmarshal the response body into the locations struct
	var locations locations
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return nil, err
	}
	// Return the names of the locations
	var locationNames []string
	for _, location := range locations.Results {
		locationNames = append(locationNames, location.Name)
	}
	return locationNames, nil
}
