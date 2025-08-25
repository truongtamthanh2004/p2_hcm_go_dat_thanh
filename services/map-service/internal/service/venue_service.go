package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"map-service/internal/config"
	"map-service/internal/dto"
	"net/http"
	"os"
	"time"
)

type VenueService interface {
	FetchVenues() ([]dto.Venue, error)
}

type venueServiceImpl struct {
	BaseURL string
}

func NewVenueService() VenueService {
	return &venueServiceImpl{
		BaseURL: os.Getenv("VENUE_SERVICE_URL"),
	}
}

func (s *venueServiceImpl) FetchVenues() ([]dto.Venue, error) {
	ctx := context.Background()
	cacheKey := "venues:all"

	if data, err := config.Get(ctx, cacheKey); err == nil {
		var venues []dto.Venue
		if json.Unmarshal([]byte(data), &venues) == nil {
			return venues, nil
		}
	}

	// resp, err := http.Get(fmt.Sprintf("%s/api/v1/venues", s.BaseURL))
	// if err != nil {
	// 	return nil, err
	// }
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/venues", s.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch venues: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch venues: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Message string      `json:"message"`
		Data    []dto.Venue `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if err := config.Set(ctx, cacheKey, string(body), 5*time.Minute); err != nil {
		log.Printf("Failed to set cache for key %s: %v", cacheKey, err)
	}

	return result.Data, nil
}
