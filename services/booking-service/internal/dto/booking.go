package dto


type CheckAvailabilityRequest struct {
	SpaceIDs  []uint    `json:"space_ids"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type CheckAvailabilityResponse struct {
	UnavailableSpaceIDs []uint `json:"unavailable_space_ids"`
}
