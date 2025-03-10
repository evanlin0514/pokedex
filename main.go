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
	"math/rand"
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

func unmarshalJson(url string, cache *pokecache.Cache, target any) error {
	data, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error getting data: %v", err)
		}
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading data: %v", err)
		}	
		cache.Add(url, data)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("error unmarshaling: %v", err)
	}
	return nil
}


func printLocate(c *config, cache *pokecache.Cache, name string) error {
	var data LocateData
	var err error
	if name == "map" {
		err = unmarshalJson(c.next, cache, &data) //data will be locateData
	} else {
		err = unmarshalJson(c.previous, cache, &data)
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

func printPokemon(cache *pokecache.Cache, area string, pokeMap *map[string]bool) error {
	var data PokeAreaData
	var err error
	(*pokeMap) = make(map[string]bool)
	if area == ""{
		return fmt.Errorf("invalid explore input: empty input")
	}

	url := "https://pokeapi.co/api/v2/location-area/" + area
	err = unmarshalJson(url, cache, &data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data: %v", err)
	}

	fmt.Printf("Exploring %v...\nFound Pokemon:\n", data.Location.Name)
	for _, pokemon := range data.PokemonEncounters{
		fmt.Printf("- %v\n", pokemon.Pokemon.Name)
		(*pokeMap)[pokemon.Pokemon.Name] = true
	}
	return nil
}

func possibilityByPercentage(length int, percentage int) bool {
    return rand.Intn(length) < percentage
}

func catchPokemon(cache *pokecache.Cache, target string, pokeMap *map[string]bool, mypokedex *map[string]PokeData) error {
	var data PokeData
	var err error
	url := "https://pokeapi.co/api/v2/pokemon/" + target

	if ok := (*pokeMap)[target]; !ok {
		return fmt.Errorf("target not in this area")
	} 
	err = unmarshalJson(url, cache, &data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data")
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", data.Forms[0].Name)
	if catch := possibilityByPercentage(data.BaseExperience, 70); catch {
		fmt.Printf("%v was caught!\n", data.Forms[0].Name)
		(*mypokedex)[data.Forms[0].Name] = data
	} else {
		fmt.Printf("%v escaped!\n", data.Forms[0].Name)
	}
	return nil
}

func inspectPokemon(target string, mypokedex *map[string]PokeData) {
	if pokemon, ok := (*mypokedex)[target]; ok {
		fmt.Printf("Name: %v\n", pokemon.Forms[0].Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats{
			fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.Base_stat)
		}
		fmt.Println("Types:")
		for _, typ := range pokemon.Types{
			fmt.Printf("  -%v\n", typ.Typ.Name)
		}
	} else {
		fmt.Println("haven't catched it yet!")
	}
	
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

type PokeAreaData struct {
	Location locateArea `json:"location"`
	PokemonEncounters []struct {
		Pokemon Pokemon `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type PokeData struct {
	BaseExperience int `json:"base_experience"`
	Weight int `json:"weight"`
	Height int `json:"height"`
	Stats []struct {
		Base_stat int `json:"base_stat"`
		Stat struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Typ struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Forms []Pokemon `json:"forms"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	page := config{
		next: "https://pokeapi.co/api/v2/location-area",
		previous: "",
	}	
	cache := pokecache.NewCache(time.Second * 5)
	pokeList := make(map[string]bool, 0)
	mypokedex := make(map[string]PokeData)
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
					"explore": {
						name: "explore",
						description: "explore the area, look up all the pokemon.",
						callback: func() error {
							return nil
						},
						page: &page,
					},
					"catch": {
						name: "catch",
						description: "catch pokemon in area you are in.",
						callback: func() error {
							return nil
						},
						page: &page,
					},
					"inspect": {
						name: "inspect",
						description: "see your pokemon's stats",
						callback: func() error {
							return nil
						},
						page: &page,
					},
				}
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			texts := scanner.Text()
			cleanSlice := cleanInput(texts)
			if len(cleanSlice) == 0 || len(cleanSlice) > 2 {
				fmt.Println("Index error: too many or empty input.")
				continue
			} 
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
					case "explore":
						if len(cleanSlice) != 2 {
							fmt.Println("error: invalid explore command input, index shoud be two.")
							continue
						}
						err := printPokemon(cache, cleanSlice[1], &pokeList)

						if err != nil {
							fmt.Println(err)
						}
					case "catch":
						if len(cleanSlice) != 2 {
							fmt.Println("error: invalid catch command input, index shoud be two.")
							continue
						}
						err := catchPokemon(cache, cleanSlice[1], &pokeList, &mypokedex)
						if err != nil {
							fmt.Println(err)
						}
					case "inspect":
						if len(cleanSlice) != 2 {
							fmt.Println("error: invalid catch command input, index shoud be two.")
							continue
						}
						inspectPokemon(cleanSlice[1], &mypokedex)
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
