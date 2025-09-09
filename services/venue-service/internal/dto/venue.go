package dto

type CreateVenueRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	City        string `json:"city"`
	Description string `json:"description"`
}

type UpdateVenueRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	City        string `json:"city"`
	Description string `json:"description"`
}

type FilterVenueRequest struct {
	Status string `form:"status"` // pending / approved / blocked
}
