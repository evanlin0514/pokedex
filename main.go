package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
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

type cliCommand struct {
	name string
	description string
	callback func() error
}


func main() {
	scanner := bufio.NewScanner(os.Stdin)	
	command := map[string]cliCommand{
					"exit": {
						name: "exit",
						description: "Exit the Pokedex.",
						callback: commandExit,
					},
					"help": {
						name: "help",
						description: "Display a help message.",
						callback: printUsage,
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
						for k, _ := range command{
							fmt.Printf("%v: %v\n", command[k].name, command[k].description)
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
