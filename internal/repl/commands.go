package repl

import (
	"fmt"
	"os"
	"text/template"

	"github.com/caleb-fringer/pokedexcli/internal/pokeapi"
)

type CommandParams any

type Handler interface {
	Execute(params CommandParams) error
}

type Command struct {
	Name        string
	Description string
	Handler
}

var registry map[string]Command
var pageState MapPagination

const helpPrompt = "Welcome to the Pokedex!\nUsage:\n\n{{range .}}{{.Name}}: {{.Description}}\n{{end}}"

func init() {
	// Initialze the value of the map's pageState
	pageState = MapPagination{
		Next:     "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		Previous: "",
	}

	// Initialize the command registry
	registry = map[string]Command{
		"exit": {
			Name:        "exit",
			Description: "Exit the pokedex",
			Handler:     ExitHandler{},
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Handler:     HelpHandler{},
		},
		"map": {
			Name:        "map",
			Description: "Get a page of location-areas",
			Handler:     MapHandler{},
		},
		"mapb": {
			Name:        "mapb",
			Description: "Get the previous page of location-areas",
			Handler:     MapBackHandler{},
		},
	}
}

// Handlers

/* Exit command
 * Takes no arguments, prints an exit message, and exits w/ code 0.
 * Always returns nil.
 */
type ExitHandler struct{}

func (h ExitHandler) Execute(args CommandParams) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

/* Help command
 * Takes no arguments and prints a help message followed by a list
 * and description of each command.
 * Panics if parsing the output template fails.
 * Returns an error if creating the output template fails.
 */
type HelpHandler struct{}

func (h HelpHandler) Execute(args CommandParams) error {
	helpTemplate := template.New("HelpTemplate")
	helpTemplate = template.Must(helpTemplate.Parse(helpPrompt))
	err := helpTemplate.Execute(os.Stdout, registry)
	if err != nil {
		return err
	}
	return nil
}

/* This struct maintains the Next and Previous links for the map and mapb
 * commands. It should ONLY be modified by MapHandler and MapBackHandler's
 * `Execute` methods.
 */
type MapPagination struct {
	Next     string
	Previous string
}

/* Map command
 * Takes no arguments. Prints the next page of map-area locations from Pokeapi.
 * Maintains the state of the current position in page results.
 * Returns an error if the handler fails to parse the provided arguments as a
 * MapPagination pointer, or if the pokeapi package returns an error.
 */
type MapHandler struct{}

func (h MapHandler) Execute(params CommandParams) error {
	response, err := pokeapi.GetLocationArea(pageState.Next)
	if err != nil {
		return err
	}

	for _, locArea := range response.Results {
		fmt.Println(locArea.Name)
	}

	pageState.Next, pageState.Previous = response.Next, response.Previous

	return nil
}

/* Mapback command
 * Takes no arguments. Prints the prev. page of map-area locations from Pokeapi,
 * or prints an message if MapBack is called while on the first page of results.
 * Returns an error if the handler fails to parse the provided arguments as a
 * MapPagination pointer, or if the pokeapi package returns an error.
 */
type MapBackHandler struct{}

func (h MapBackHandler) Execute(params CommandParams) error {
	if pageState.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	response, err := pokeapi.GetLocationArea(pageState.Previous)
	if err != nil {
		return err
	}

	for _, locArea := range response.Results {
		fmt.Println(locArea.Name)
	}

	pageState.Next, pageState.Previous = response.Next, response.Previous
	return nil
}
