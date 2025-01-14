package user

// CreateRequest represents the request payload for creating a user.
// swagger:model CreateRequest
type CreateRequest struct {
	Name        string  `json:"name" binding:"required,min=3"`
	DateOfBirth string  `json:"date_of_birth" binding:"required,datetime=2006-01-02"`
	Address     string  `json:"address" binding:"required"`
	Description *string `json:"description,omitempty"`
	Latitude    float64 `json:"latitude" binding:"required,latitude"`
	Longitude   float64 `json:"longitude" binding:"required,longitude"`
}

// CreateResp represents the successful response when a user is created.
// swagger:model CreateResp
type CreateResp struct {
	ID string `json:"id"`
}

// GetResp represents the successful response when a user is retrieved.
// swagger:model GetResp
type GetResp struct {
	Name        string  `json:"name"`
	DateOfBirth string  `json:"date_of_birth"`
	Address     string  `json:"address"`
	Description *string `json:"description,omitempty"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
