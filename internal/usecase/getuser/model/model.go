package model

import "time"

type User struct {
	ID          string
	Name        string
	DateOfBirth string
	Address     string
	Location    Point
	Description *string
	CreatedAt   time.Time
}

type Point struct {
	Longitude float64
	Latitude  float64
}
