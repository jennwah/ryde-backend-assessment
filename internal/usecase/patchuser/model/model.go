package model

type PatchUser struct {
	ID          string
	Name        string
	DateOfBirth string
	Address     string
	Latitude    float64
	Longitude   float64
	Description *string
}
