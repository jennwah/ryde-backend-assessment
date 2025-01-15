package model

type Friend struct {
	ID          string
	Name        string
	DateOfBirth string
	Address     string
	Location    Point
	Description *string
}

type Point struct {
	Longitude float64
	Latitude  float64
}

type Friends []Friend
