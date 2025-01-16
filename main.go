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
	Callback    func() error
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
	}

	tokenizer = regexp.MustCompile("[[:alpha:]]+")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if ok := scanner.Scan(); !ok {
			fmt.Println()
			os.Exit(0)
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
	commandStruct, ok := cmdRegistry[command]
	if !ok {
		return false
	}

	err := commandStruct.Callback()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// Callbacks for commands
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	helpTemplate := template.New("HelpTemplate")
	helpTemplate = template.Must(helpTemplate.Parse("Welcome to the Pokedex!\nUsage:\n\n{{range .}}{{.Name}}: {{.Description}}\n{{end}}"))
	err := helpTemplate.Execute(os.Stdout, cmdRegistry)
	if err != nil {
		return err
	}
	return nil
}
