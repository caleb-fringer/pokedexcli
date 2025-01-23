package pokeapi

import "testing"

func TestGetLocationArea(t *testing.T) {
	response, err := GetLocationAreas(0, 20)
	if err != nil {
		t.Fatalf("Querying https://pokeapi.co/api/v2/location-area?offset=0&limit=20 returned an error: %v", err)
	}

	firstArea := response.Results[0].Name
	if firstArea != "canalave-city-area" {
		t.Fatalf("Querying the first location-area returned wrong area.\n\tExpected: %s\n\tFound: %s", firstArea, "canalave-city-area")
	}
}
