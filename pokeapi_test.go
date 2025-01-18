package main

import "testing"

func TestGetLocationArea(t *testing.T) {
	endpoint := pageLink{
		"https://pokeapi.co/api/v2/location-area",
		"",
	}

	results, err := getLocationArea(&endpoint)
	if err != nil {
		t.Fatalf("Querying %s returned an error", endpoint.Next)
	}

	if results[0].Name != "canalave-city-area" {
		t.Fatalf("Querying the first location-area returned wrong area.\n\tExpected: %s\n\tFound: %s", results[0].Name, "canalave-city-area")
	}

	if endpoint.Next != "https://pokeapi.co/api/v2/location-area?offset=20&limit=20" {
		t.Fatalf("Failed to update the pageLink structure correctly.\n\tExpected: %s\n\tFound: %s",
			"https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
			endpoint.Next)
	}
}
