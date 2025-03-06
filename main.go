package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string{
	lowerText := strings.ToLower(text)
	cleanedList := strings.Split(lowerText, " ")
	return cleanedList
}

func main() {
	fmt.Println(cleanInput("Charmander Bulbasaur PIKACHU"))
}
