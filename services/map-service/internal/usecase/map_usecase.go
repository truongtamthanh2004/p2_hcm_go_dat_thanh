package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"map-service/internal/config"
	"map-service/internal/dto"
	"map-service/internal/service"
	"time"
)

type MapUsecase struct {
	VenueService service.VenueService
}

func NewMapUsecase(vs service.VenueService) *MapUsecase {
	return &MapUsecase{VenueService: vs}
}

func (uc *MapUsecase) GetVenues() ([]dto.Venue, error) {
	return uc.VenueService.FetchVenues()
}

func (u *MapUsecase) GetVenuesWithLocation() ([]dto.Venue, error) {
	ctx := context.Background()

	venues, err := u.VenueService.FetchVenues()
	if err != nil {
		return nil, err
	}

	for i, v := range venues {
		coord, err := service.GeocodeAddress(fmt.Sprintf("%s, %s", v.Address, v.City))
		if err != nil {
			log.Printf("⚠️ Geocode failed for venue %q (%s, %s): %v", v.Name, v.Address, v.City, err)
			continue
		}
		venues[i].Latitude = coord.Lat
		venues[i].Longitude = coord.Lng
	} 

	jsonStr, err := toJSON(venues)
	if err != nil {
		log.Printf("❌ Failed to marshal venues for cache: %v", err)
	} else {
		if err := config.Set(ctx, "venues:all", jsonStr, 5*time.Minute); err != nil {
			log.Printf("❌ Failed to set venues cache: %v", err)
		}
	}

	return venues, nil
}

func toJSON(v interface{}) (string, error) {
    b, err := json.Marshal(v)
    if err != nil {
        return "", err
    }
    return string(b), nil
}
