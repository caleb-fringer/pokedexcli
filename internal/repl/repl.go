package repl

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/caleb-fringer/pokedexcli/internal/pokeapi"
)

var tokenizer *regexp.Regexp
var exploreArgValidator *regexp.Regexp

func init() {
	//tokenizer = regexp.MustCompile("[[:alpha:]]+")
	tokenizer = regexp.MustCompile("[[:alpha:]]+(?:-[[:alnum:]]+)*")
}

func DoREPL() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if ok := scanner.Scan(); !ok {
			fmt.Println()
			os.Exit(1)
		}

		line := scanner.Text()
		tokens := cleanInput(line)

		if len(tokens) == 0 {
			continue
		}

		cmd := tokens[0]
		args := tokens[1:]
		doCommand(cmd, args)
	}
}

func cleanInput(text string) (tokens []string) {
	lower := strings.ToLower(text)
	return tokenizer.FindAllString(lower, -1)
}

func doCommand(command string, args []string) bool {
	// Fetch the command structure, returning if not found.
	commandStruct, ok := registry[command]
	if !ok {
		fmt.Println("Please provide a supported command. Try `help` if you don't know them!")
		return false
	}

	// Populate the correct CommandParams struct according to the cmd called.
	var params CommandParams

	switch command {
	// Here we will populate special CommandParam structs as needed.
	case "explore":
		if len(args) < 1 {
			fmt.Println("Please provide a location-area to explore!")
			return false
		}
		params = args[0]
	case "catch", "inspect":
		if len(args) < 1 {
			fmt.Println("Please provide a Pokemon to capture!")
			return false
		}
		params = args[0]
	}

	err := commandStruct.Execute(params)
	if err != nil {
		// Ignore ResourceNotFoundErrors, they do not need to be handled.
		if _, ok := err.(pokeapi.ResourceNotFoundError); !ok {
			fmt.Println(err)
		}
		return false
	}
	return true
}
