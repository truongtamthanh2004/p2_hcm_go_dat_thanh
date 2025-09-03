package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type VenueService interface {
	GetSpaceByID(spaceID uint) (*Space, error)
}

type venueHTTPService struct {
	baseURL string
	client  *http.Client
}

type spaceResponse struct {
  Data    Space  `json:"data"`
  Message string `json:"message"`
}

type Space struct {
  ID          uint    `json:"ID"`
  VenueID     uint    `json:"VenueID"`
  Name        string  `json:"Name"`
  Capacity    int     `json:"Capacity"`
  Price       float64 `json:"Price"`
  Description string  `json:"Description"`
}

func NewVenueHTTPService(baseURL string) VenueService {
	return &venueHTTPService{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *venueHTTPService) GetSpaceByID(spaceID uint) (*Space, error) {
  url := fmt.Sprintf("%s/api/v1/spaces/%d", s.baseURL, spaceID)
  resp, err := s.client.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("venue service returned status %d", resp.StatusCode)
  }

  var result spaceResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
    }

  return &result.Data, nil
}
