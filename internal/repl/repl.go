package repl

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var tokenizer *regexp.Regexp

func init() {
	tokenizer = regexp.MustCompile("[[:alpha:]]+")
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
		doCommand(cmd)
	}
}

func cleanInput(text string) (tokens []string) {
	lower := strings.ToLower(text)
	return tokenizer.FindAllString(lower, -1)
}

func doCommand(command string) bool {
	// Fetch the command structure, returning if not found.
	commandStruct, ok := registry[command]
	if !ok {
		return false
	}

	// Populate the correct CommandParams struct according to the cmd called.
	var params CommandParams
	switch command {
	// Here we will populate special CommandParam structs as needed.
	}

	err := commandStruct.Execute(params)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
