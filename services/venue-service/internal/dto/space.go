package dto

type CreateSpaceRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required,oneof=private_office meeting_room desk"`
	Capacity    int     `json:"capacity" binding:"required,min=1"`
	Price       float64 `json:"price" binding:"required"`
	Description string  `json:"description"`
	OpenHour    string  `json:"open_hour" binding:"required"`
	CloseHour   string  `json:"close_hour" binding:"required"`
}

type UpdateSpaceRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type" binding:"omitempty,oneof=private_office meeting_room desk"`
	Capacity    int     `json:"capacity" binding:"omitempty,min=1"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	OpenHour    string  `json:"open_hour"`
	CloseHour   string  `json:"close_hour"`
}

type UpdateManagerRequest struct {
	ManagerID uint `json:"manager_id" binding:"required"`
}
