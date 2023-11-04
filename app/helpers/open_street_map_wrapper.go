package helpers

import (
    "fmt"
    "net/http"
    "io/ioutil"
		"encoding/json"
		"strconv"
		"net/url"
)

type AddressResult struct {
	PlaceID      int     `json:"place_id"`
	Licence      string  `json:"licence"`
	OSMType      string  `json:"osm_type"`
	OSMID        int     `json:"osm_id"`
	Latitude     string  `json:"lat"`
	Longitude    string  `json:"lon"`
	Class        string  `json:"class"`
	Type         string  `json:"type"`
	PlaceRank    int     `json:"place_rank"`
	Importance   float64 `json:"importance"`
	AddressType  string  `json:"addresstype"`
	Name         string  `json:"name"`
	DisplayName  string  `json:"display_name"`
	BoundingBox   []string `json:"boundingbox"`
}

func RetrieveCoordinatesFromAddress(address string) (float64, float64, error){
    // address = "1600 Amphitheatre Parkway, Mountain View, CA"
    // Create the URL for the Nominatim API request
		apiUrl := "https://nominatim.openstreetmap.org/search"
		encodedQuery := url.QueryEscape(address)
    apiUrl += "?q=" + encodedQuery
		apiUrl += "&format=json"

		// Create a new HTTP request
    req, err := http.NewRequest("GET", apiUrl, nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return 0, 0, err
    }

    // Set the User-Agent header
		req.Header.Set("User-Agent", "whir/1.0")

    // Send the request
    client := &http.Client{}
    response, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return 0, 0, err
    }
    defer response.Body.Close()

    // Read and parse the JSON response
    responseBody, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        return 0, 0, err
		}

    // Unmarshal the JSON data to a slice of results
		var results []AddressResult
    if err := json.Unmarshal(responseBody, &results); err != nil {
        fmt.Println("Error parsing JSON:", err)
        return 0, 0, err
		}
		lat := 0.0
		lon := 0.0

    // Extract the coordinates from the first result
    if len(results) > 0 {
			    // Convert the string to a float64
					lat, err = strconv.ParseFloat(results[0].Latitude, 64)
					if err != nil {
							fmt.Println("Error:", err)
							return 0, 0, err
					}
					lon, err = strconv.ParseFloat(results[0].Longitude, 64)
					if err != nil {
							fmt.Println("Error:", err)
							return 0, 0, err
					}
    } else {
        fmt.Println("No results found")
		}
		return lon, lat, nil
}