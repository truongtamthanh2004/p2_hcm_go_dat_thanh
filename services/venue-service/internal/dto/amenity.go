package dto

type AddAmenityRequest struct {
	AmenityID uint `json:"amenity_id" binding:"required"`
}

type CreateAmenityRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

type UpdateAmenityRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=255"`
}