package pokeapi

import "testing"

func TestGetLocationAreas(t *testing.T) {
	response, err := GetLocationAreas(0, 20)
	if err != nil {
		t.Fatalf("Querying https://pokeapi.co/api/v2/location-area?offset=0&limit=20 returned an error: %v", err)
	}

	firstArea := response.Results[0].Name
	if firstArea != "canalave-city-area" {
		t.Fatalf("Querying the first location-area returned wrong area.\n\tExpected: %s\n\tFound: %s", firstArea, "canalave-city-area")
	}
}

func TestGetLocationArea(t *testing.T) {
	response, err := GetLocationArea("pastoria-city-area")
	if err != nil {
		t.Fatalf("Querying https://pokeapi.co/api/v2/location-area?offset=0&limit=20 returned an error: %v", err)
	}

	firstPokemon := response.PokemonEncounters[0].Pokemon.Name
	if firstPokemon != "tentacool" {
		t.Fatalf("Querying the pastoria-city-area returned wrong first Pokemon.\n\tExpected: %s\n\tFound: %s",
			firstPokemon,
			"tentacool")
	}
}
