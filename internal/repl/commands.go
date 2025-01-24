package repl

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
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

const pageSize = 20
const helpPrompt = "Welcome to the Pokedex!\nUsage:\n\n{{range .}}{{.Name}}: {{.Description}}\n{{end}}"

func init() {
	// Initialze the value of the map's pageState
	initial_url, err := url.Parse(
		fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=0&limit=%d",
			pageSize))
	if err != nil {
		log.Fatal(err)
	}

	pageState = MapPagination{
		Next:     initial_url,
		Previous: &url.URL{},
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
		"explore": {
			Name:        "explore",
			Description: "Explore a location-area for Pokemon",
			Handler:     ExploreHandler{},
		},
		"catch": {
			Name:        "catch",
			Description: "Catch the given Pokemon",
			Handler:     CatchHandler{},
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
	Next     *url.URL
	Previous *url.URL
}

/* Map command
 * Takes no arguments. Prints the next page of map-area locations from Pokeapi.
 * Maintains the state of the current position in page results.
 * Returns an error if the handler fails to extract the page offset and limit
 * from the pageState.Next and pageState.Previous urls, or if the pokeapi
 * package returns an error.
 */
type MapHandler struct{}

func (h MapHandler) Execute(params CommandParams) error {
	queryParams := pageState.Next.Query()

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil {
		return fmt.Errorf("Error parsing offset query param to int: %w", err)
	}

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		return fmt.Errorf("Error parsing limit query param to int: %w", err)
	}

	response, err := pokeapi.GetLocationAreas(offset, limit)
	if err != nil {
		return err
	}

	for _, locArea := range response.Results {
		fmt.Println(locArea.Name)
	}
	fmt.Println()

	err = pageState.updateState(response.Next, response.Previous)
	if err != nil {
		return fmt.Errorf("Error updating MapPagination state: %w", err)
	}

	return nil
}

/* Mapback command
 * Takes no arguments. Prints the prev. page of map-area locations from Pokeapi,
 * or prints an message if MapBack is called while on the first page of results.
 * Returns an error if the handler fails to extract the page offset and limit
 * from the pageState.Next and pageState.Previous urls, or if the pokeapi
 * package returns an error.
 */
type MapBackHandler struct{}

func (h MapBackHandler) Execute(params CommandParams) error {
	if pageState.Previous.Path == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	queryParams := pageState.Previous.Query()

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil {
		return fmt.Errorf("Error parsing offset query param to int: %w", err)
	}

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		return fmt.Errorf("Error parsing limit query param to int: %w", err)
	}

	response, err := pokeapi.GetLocationAreas(offset, limit)
	if err != nil {
		return err
	}

	for _, locArea := range response.Results {
		fmt.Println(locArea.Name)
	}
	fmt.Println()

	err = pageState.updateState(response.Next, response.Previous)

	if err != nil {
		return fmt.Errorf("Error updating MapPagination state: %w", err)
	}

	return nil
}

func (m *MapPagination) updateState(next, prev string) error {
	newNext, err := url.Parse(next)
	if err != nil {
		return fmt.Errorf("Error updating MapPagination.Next: %w", err)
	}

	newPrev, err := url.Parse(prev)
	if err != nil {
		return fmt.Errorf("Error updating MapPagination.Prev: %w", err)
	}

	m.Next, m.Previous = newNext, newPrev

	return nil
}

/* Explore command.
 * Takes the name of a location-area to explore, and prints a list of all
 * Pokemon at that location, or "Location not found" if the pokeapi returns a
 * status code 404.
 *
 * Returns an error if the handler fails to coerce the provided arguments as a
 * string, or if the pokeapi package returns an error.
 */
type ExploreHandler struct{}

func (h ExploreHandler) Execute(params CommandParams) error {
	locationAreaName, ok := params.(string)
	if !ok {
		return errors.New("Failed type assertion to string. ExploreHandler requires a string argument")
	}

	fmt.Printf("Exploring %v...\n", locationAreaName)

	response, err := pokeapi.GetLocationArea(locationAreaName)
	if err != nil {
		switch err.(type) {
		case pokeapi.ResourceNotFoundError:
			fmt.Println("Location not found!")
			return err
		default:
			return fmt.Errorf("Error fetching requested location-area: %w", err)
		}
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range response.PokemonEncounters {
		fmt.Printf("\t- %s\n", pokemon.Pokemon.Name)
	}
	fmt.Println()

	return nil
}

/* Catch command
 * TODO: Document me!
 */
type CatchHandler struct{}

func (h CatchHandler) Execute(params CommandParams) error {
	pokemonName, ok := params.(string)
	if !ok {
		return errors.New("Failed type assertion to string. CatchHandler requires a string argument")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	response, err := pokeapi.GetPokemon(pokemonName)
	if err != nil {
		switch err.(type) {
		case pokeapi.ResourceNotFoundError:
			fmt.Println("Pokemon not found!")
			return err
		default:
			return fmt.Errorf("Error fetching requested Pokemon: %w", err)
		}
	}

	fmt.Printf("Base experience: %d\n", response.BaseExperience)

	return nil
}
