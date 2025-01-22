package pokeapi

import "testing"

func TestGetLocationArea(t *testing.T) {
	endpoint := "https://pokeapi.co/api/v2/location-area"

	response, err := GetLocationArea(endpoint)
	if err != nil {
		t.Fatalf("Querying %s returned an error", endpoint)
	}

	firstArea := response.Results[0].Name
	if firstArea != "canalave-city-area" {
		t.Fatalf("Querying the first location-area returned wrong area.\n\tExpected: %s\n\tFound: %s", firstArea, "canalave-city-area")
	}
}
