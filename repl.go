package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type cliCommand struct {
	Name        string
	Description string
	Callback    func(c *pageLink) error
}

var cmdRegistry map[string]cliCommand
var tokenizer *regexp.Regexp

func init() {
	cmdRegistry = map[string]cliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the pokedex",
			Callback:    commandExit,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Get a page of location-areas",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "map-back",
			Description: "Get the previous page of location-areas",
			Callback:    commandMapback,
		},
	}

	tokenizer = regexp.MustCompile("[[:alpha:]]+")
}

func doREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	var c pageLink

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
		doCommand(cmd, &c)
	}
}

func cleanInput(text string) (tokens []string) {
	lower := strings.ToLower(text)
	return tokenizer.FindAllString(lower, -1)
}

func doCommand(command string, c *pageLink) bool {
	commandStruct, ok := cmdRegistry[command]
	if !ok {
		return false
	}

	err := commandStruct.Callback(c)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// Callbacks for commands. Each command will return an optional error
func commandExit(c *pageLink) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *pageLink) error {
	helpTemplate := template.New("HelpTemplate")
	helpTemplate = template.Must(helpTemplate.Parse(`Welcome to the Pokedex!
Usage:

{{range .}}{{.Name}}: {{.Description}}
{{end}}`))
	err := helpTemplate.Execute(os.Stdout, cmdRegistry)
	if err != nil {
		return err
	}
	return nil
}

func commandMap(c *pageLink) error {
	return nil
}

func commandMapback(c *pageLink) error {
	return nil
}
