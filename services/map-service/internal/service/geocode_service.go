package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type GeoCodingResult struct {
	Lat float64
	Lng float64
}

func GeocodeAddress(address string) (GeoCodingResult, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
    return GeoCodingResult{}, fmt.Errorf("missing GOOGLE_MAPS_API_KEY environment variable")
  }

	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"

	u := fmt.Sprintf("%s?address=%s&key=%s", baseURL, url.QueryEscape(address), apiKey)
	resp, err := http.Get(u)
	if err != nil {
		return GeoCodingResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(resp.Body)
    return GeoCodingResult{}, fmt.Errorf("google geocoding API returned %d: %s", resp.StatusCode, string(body))
  }

	var result struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GeoCodingResult{}, err
	}

	if result.Status != "OK" || len(result.Results) == 0 {
		return GeoCodingResult{}, fmt.Errorf("geocode failed for address: %s", address)
	}

	return GeoCodingResult{
		Lat: result.Results[0].Geometry.Location.Lat,
		Lng: result.Results[0].Geometry.Location.Lng,
	}, nil
}
