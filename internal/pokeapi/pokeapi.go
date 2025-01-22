package pokeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/caleb-fringer/pokedexcli/internal/pokecache"
)

var cache pokecache.Cache

func init() {
	cache = pokecache.NewCache(5 * time.Second)
}

/* GetLocationArea
 * Given a url,  this function will:
 *     -GETs the url from PokeAPI
 *     -Cache the raw data
 *     -Unmarshal the JSON response into a LocationAreaResponse object
 *     -Return the response.
 *
 * If the url is already in  the cache, it will unmarshal that data source
 * instead.
 *
 * An error occurs if the url is blank, if the http.GET call fails, if the
 * response's status code is not 200, or if decoding the response fails.
 */
func GetLocationArea(url string) (response LocationAreaResponse, err error) {
	if url == "" {
		return response, fmt.Errorf("Uninitialized url")
	}

	var raw []byte
	// Check cache for data. Store the raw data for unmarshalling
	if data, ok := cache.Get(url); ok {
		raw, err = io.ReadAll(bytes.NewReader(data))
		if err != nil {
			return response, fmt.Errorf("Error reading cached raw data: %w", err)
		}
	} else {
		// Make HTTP request and cache result
		res, err := http.Get(url)
		if err != nil {
			return response, fmt.Errorf("HTTP error when GET'ing %s: %w", url, err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return response, fmt.Errorf("Invalid HTTP response code from %s, status: %d", url, res.StatusCode)
		}

		// Cache the result for later use
		raw, err = io.ReadAll(res.Body)
		if err != nil {
			return response, fmt.Errorf("Error reading raw response data to add to cache: %w", err)
		}
		cache.Add(url, raw)
	}

	// Unmarshall response
	err = json.Unmarshal(raw, &response)

	if err != nil {
		return response, fmt.Errorf("Error unmarshalling response from %s: %w", url, err)
	}
	return response, nil
}
