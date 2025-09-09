package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BookingClient interface {
	CheckAvailability(ctx context.Context, spaceIDs []uint, start, end time.Time) ([]uint, error)
}

type bookingClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewBookingClient(baseURL string) BookingClient {
	return &bookingClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type checkAvailabilityRequest struct {
	SpaceIDs  []uint `json:"space_ids"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type checkAvailabilityResponse struct {
	UnavailableSpaceIDs []uint `json:"unavailable_space_ids"`
}

func (c *bookingClient) CheckAvailability(ctx context.Context, spaceIDs []uint, start, end time.Time) ([]uint, error) {
	reqBody := checkAvailabilityRequest{
		SpaceIDs:  spaceIDs,
		StartTime: start.UTC().Format(time.RFC3339),
		EndTime:   end.UTC().Format(time.RFC3339),
	}
	body, err := json.Marshal(reqBody)
	fmt.Println("DEBUG request to booking-service:", string(body))
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/internal/bookings/check-availability", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("booking-service returned status %d", resp.StatusCode)
	}

	var res checkAvailabilityResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.UnavailableSpaceIDs, nil
}
