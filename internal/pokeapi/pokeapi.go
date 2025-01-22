/* This package is responsible for making HTTP GET requests to pokeapi.co
 * endpoints, caching responses, and unmarshalling JSON responses. All further
 * data processing, extracting fields, and page management should be handled
 * by the specific command handlers that consume this API.
 */
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/caleb-fringer/pokedexcli/internal/pokecache"
)

var cache pokecache.Cache

func init() {
	cache = pokecache.NewCache(5 * time.Second)
}

/* GetLocationAreas
 * Given a url, this function will:
 *     -GETs the url from PokeAPI
 *     -Cache the raw data
 *     -Unmarshal the JSON response into a LocationAreaResponse object
 *     -Return the response.
 *
 * If the url is already in  the cache, it will unmarshal that data source
 * instead.
 *
 * Returns an error occurs if the url is blank, if the http.GET call fails, if
 * the response's status code is not 200, or if decoding the response fails.
 */
func GetLocationAreas(url url.URL) (response LocationAreasResponse, err error) {
	if url.Path == "" {
		return response, fmt.Errorf("Uninitialized url")
	}

	data, ok := isCached(url)

	// Make HTTP request and cache result on cache miss
	if !ok {
		res, err := http.Get(url.String())
		if err != nil {
			return response, fmt.Errorf("HTTP error when GET'ing %s: %w", url, err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return response, fmt.Errorf("Invalid HTTP response code from %s, status: %d", url, res.StatusCode)
		}

		// Cache the result for later use
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return response, fmt.Errorf("Error reading raw response data to add to cache: %w", err)
		}
		cache.Add(url, data)
	}

	// Unmarshall response
	err = json.Unmarshal(data, &response)

	if err != nil {
		return response, fmt.Errorf("Error unmarshalling response from %s: %w", url, err)
	}
	return response, nil
}

/* GetLocationArea
 * Given a specific LocationArea name, this function will:
 *     -GETs the url from PokeAPI
 *     -Cache the raw data
 *     -Unmarshal the JSON response into a LocationAreaResponse object
 *     -Return the response.
 *
 * If the url is already in  the cache, it will unmarshal that data source
 * instead.
 *
 * Returns an error occurs if the url is blank, if the http.GET call fails, if
 * the response's status code is not 200, or if decoding the response fails.
 */
func GetLocationArea(name string) (response LocationAreaResponse, err error) {
	return
}

/* isCached:
 * Helper function that checks cache for a requested resource, returning the
 * raw data and a boolean indicating if the cache hit or not
 */
func isCached(url url.URL) (data []byte, ok bool) {
	data, ok = cache.Get(url)
	if !ok {
		return nil, false
	}
	return
}
