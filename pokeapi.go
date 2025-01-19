package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/caleb-fringer/pokedexcli/internal/pokecache"
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

var cache pokecache.Cache

func init() {
	cache = pokecache.NewCache(5 * time.Second)
}

// Given a pageLink struct, this function GETS the pages.Next url from PokeAPI,
// caches the raw data, unmarshals the JSON response into a slice of
// LocationArea objects, and returns the slice. If the url is already in  the
// cache, it will unmarshal that data source instead.
// An error occurs if the pages struct does not have an pages.Next string,
// if the http.GET call fails, if the response's status code is not 200,
// or if decoding the response fails.
func getLocationArea(pages *pageLink) (results []LocationArea, err error) {
	if pages.Next == "" {
		return nil, fmt.Errorf("Uninitialized pages struct.")
	}

	var raw []byte
	// Check cache for data. Store the raw data for unmarshalling
	if data, ok := cache.Get(pages.Next); ok {
		raw, err = io.ReadAll(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("Error reading cached raw data: %w", err)
		}
	} else {
		// Make HTTP request and cache result
		res, err := http.Get(pages.Next)
		if err != nil {
			return nil, fmt.Errorf("HTTP error when GET'ing %s: %w", pages.Next, err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, fmt.Errorf("Invalid HTTP response code from %s, status: %d", pages.Next, res.StatusCode)
		}

		// Cache the result for later use
		raw, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Error reading raw response data to add to cache: %w", err)
		}
		cache.Add(pages.Next, raw)
	}

	// Unmarshall response and extract results
	var response LocationAreaResponse
	err = json.Unmarshal(raw, &response)

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling response from %s: %w", pages.Next, err)
	}

	results = response.Results
	pages.Next = response.Next
	pages.Previous = response.Previous
	return results, nil
}
