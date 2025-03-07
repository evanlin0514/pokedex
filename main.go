package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/evanlin0514/pokedexcli/internal/pokecache"
)

func cleanInput(text string) []string{
	trimText := strings.TrimSpace(text)
	lowerText := strings.ToLower(trimText)
	cleanedList := strings.Fields(lowerText)
	return cleanedList
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func printUsage() error {
	fmt.Println("Welcome to the Pokedex!\nUsage: ")
	return nil
}

func unmarshalJson(url string, cache *pokecache.Cache) (LocateData ,error) {
	var result LocateData
	data, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return LocateData{}, fmt.Errorf("error getting data: %v", err)
		}
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return LocateData{}, fmt.Errorf("error reading data: %v", err)
		}	
		cache.Add(url, data)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("error unmarshaling: %v", err)
	}
	return result, nil
}


func printLocate(c *config, cache *pokecache.Cache, name string) error {
	var data LocateData
	var err error
	if name == "map" {
		data, err = unmarshalJson(c.next, cache) //data will be locateData	
	} else {
		data, err = unmarshalJson(c.previous, cache) //data will be locateData
	}
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	c.next = data.Next
	c.previous = data.Previous

	// check if it's first or last page
	if c.previous == ""{
		fmt.Println("--------------------")
		fmt.Println("It's the first page!")
		fmt.Println("--------------------")
	} 
	if c.next == ""{
		fmt.Println("--------------------")
		fmt.Println("It's the last page!")
		fmt.Println("--------------------")
	} 

	for _, location := range data.Results{ // location will be a single locateData struct, loop for 20
		fmt.Printf("%v\n", location.Name )
	}
	return nil
}

type cliCommand struct {
	name string
	description string
	callback func() error
	page *config
}

type config struct{
	next string
	previous string
}

type locateArea struct {
	Name string `json:"name"`
}

type LocateData struct {
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []locateArea `json:"results"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	page := config{
		next: "https://pokeapi.co/api/v2/location-area",
		previous: "",
	}	
	cache := pokecache.NewCache(time.Second * 5)
	command := map[string]cliCommand{
					"exit": {
						name: "exit",
						description: "Exit the Pokedex.",
						callback: commandExit,
						page: &page,
					},
					"help": {
						name: "help",
						description: "Display a help message.",
						callback: printUsage,
						page: &page,
					},
					"map": {
						name: "map",
						description: "Display next 20 locations.",
						callback: func() error {
							return printLocate(&page, cache, "map")
						},	
						page: &page,					
					},
					"mapb": {
						name: "mapb",
						description: "Display previous 20 locations.",
						callback: func() error {
							return printLocate(&page, cache, "mapb")
						},	
						page: &page,					
					},
				}
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			texts := scanner.Text()
			if texts == "" {
				fmt.Println("Empty input detected, try again")
				continue
			} 
			cleanSlice := cleanInput(texts)
			if key, ok := command[cleanSlice[0]]; ok{
				switch key.name {
					case  "exit":
						err := key.callback()
						if err != nil {
							fmt.Println("Error exiting Pokedex")
						}
					case  "help":
						printUsage()
						for k := range command{
							fmt.Printf("%v: %v\n", command[k].name, command[k].description)
						}
					case "map":
						err := key.callback()
						if err != nil {
							fmt.Println(err)
						}
					case "mapb":
						err := key.callback()
						if err != nil {
							fmt.Println(err)
						}
				}
			} else {
				fmt.Println("Unknown Command")
			}

		}
		
		if err := scanner.Err(); err != nil {
			fmt.Println("error scanning input.")
		}
	}
}
