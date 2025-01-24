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
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/caleb-fringer/pokedexcli/internal/pokecache"
)

var BaseUrl *url.URL
var cache *pokecache.Cache

func init() {
	var err error
	BaseUrl, err = url.Parse("https://pokeapi.co/api/v2/")
	if err != nil {
		log.Fatal("Error parsing base URL for pokeapi: %w", err)
	}
	cache = pokecache.NewCache(5 * time.Second)
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

/* GetLocationAreas
 * Given a page offset and limit:
 *     -Construct the url for the requested endpoint
 *     -GETs the resource from PokeAPI
 *     -Cache the raw data
 *     -Unmarshal the JSON response into a LocationAreasResponse object
 *     -Return the response.
 *
 * If the url is already in  the cache, it will unmarshal that data source
 * instead.
 *
 * Default values for offset, limit should be 0, 20 to request a single page of
 * 20 LocationAreas
 * Returns an error occurs if the url is blank, if the http.GET call fails, if
 * the response's status code is not 200, or if decoding the response fails.
 */
func GetLocationAreas(offset, limit int) (response LocationAreasResponse, err error) {
	// Construct query params
	queryParams := url.Values{}
	queryParams.Add("offset", strconv.Itoa(offset))
	queryParams.Add("limit", strconv.Itoa(limit))

	// Construct url w/ populated query params
	url := BaseUrl.JoinPath("location-area")
	url.RawQuery = queryParams.Encode()
	data, ok := isCached(*url)

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
		cache.Add(*url, data)
	}

	// Unmarshall response
	err = json.Unmarshal(data, &response)

	if err != nil {
		return response, fmt.Errorf("Error unmarshalling response from %s: %w", url, err)
	}
	return response, nil
}

// This is a special error to indicate that a 404 error occured so the caller
// of GetLocationArea may distinguish between a bad resource name and other
// more problematic errors.
type LocationNotFoundError struct {
	StatusCode   int
	LocationArea string
}

func (e LocationNotFoundError) Error() string {
	return fmt.Sprintf("Location-area %v not found: Status code %d", e.LocationArea, e.StatusCode)
}

/* GetLocationArea
 * Given a specific LocationArea name, this function will:
 *     -Construct the url for the requested endpoint
 *     -GETs the resource from PokeAPI
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
	// Construct the url for the requested resource
	url := BaseUrl.JoinPath("location-area", name)

	// Check if the resource is cached
	data, ok := isCached(*url)

	// Make HTTP request and cache result on cache miss
	if !ok {
		res, err := http.Get(url.String())
		if err != nil {
			return response, fmt.Errorf("HTTP error when GET'ing %s: %w", url, err)
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			return response, LocationNotFoundError{res.StatusCode, name}
		}
		if res.StatusCode != http.StatusOK {
			return response, fmt.Errorf("Invalid HTTP response code from %s, status: %d", url, res.StatusCode)
		}

		// Cache the result for later use
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return response, fmt.Errorf("Error reading raw response data to add to cache: %w", err)
		}
		cache.Add(*url, data)
	}

	// Unmarshall response
	err = json.Unmarshal(data, &response)

	if err != nil {
		return response, fmt.Errorf("Error unmarshalling response from %s: %w", url, err)
	}
	return response, nil
}
