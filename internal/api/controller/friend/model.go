package friend

// CreateReq represents the request payload for creating a user friend pair.
// swagger:model CreateReq
type CreateReq struct {
	FriendID string `json:"friend_id" binding:"required,uuid"`
}

// Header represents the header definition for friend controller.
// swagger:model Header
type Header struct {
	UserId string `header:"user_id" binding:"required,uuid"`
}

// GetNearbyReq represents the request payload for finding user's nearby friends in radius meter specified.
// swagger:model GetNearbyReq
type GetNearbyReq struct {
	RadiusMeter int `json:"radius_meter" binding:"required"`
}

// GetNearbyResp represents the response payload for finding user's nearby friends in radius meter specified.
// swagger:model GetNearbyResp
type GetNearbyResp struct {
	Friends []Friend `json:"friends" binding:"required"`
}

// Friend represents the user model.
// swagger:model Friend
type Friend struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	DateOfBirth string  `json:"date_of_birth"`
	Address     string  `json:"address"`
	Description *string `json:"description,omitempty"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
