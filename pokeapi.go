package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type pageLink struct {
	Next     string
	Previous string
}

type LocationArea struct {
	Name string
	Url  string
}

type LocationAreaResponse struct {
	pageLink
	Count   int
	Results []LocationArea
}

func getLocationArea(pages *pageLink) (results []LocationArea, err error) {
	if pages.Next == "" {
		return nil, fmt.Errorf("Uninitialized pages struct.")
	}

	res, err := http.Get(pages.Next)
	if err != nil {
		return nil, fmt.Errorf("HTTP error when GET'ing %s: %w", pages.Next, err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid HTTP response code from %s, status: %d", pages.Next, res.StatusCode)
	}

	var response LocationAreaResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		res.Body.Close()
		return nil, fmt.Errorf("Error unmarshalling response from %s: %w", pages.Next, err)
	}

	results = response.Results
	pages.Next = response.Next
	pages.Previous = response.Previous
	return results, nil
}
