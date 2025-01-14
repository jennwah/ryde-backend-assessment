package repository

type dbUser struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	DateOfBirth string  `db:"date_of_birth"`
	Address     string  `db:"address"`
	Description *string `db:"description"`
	Latitude    float64 `db:"latitude"`
	Longitude   float64 `db:"longitude"`
}
