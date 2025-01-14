package repository

import "github.com/jennwah/ryde-backend-engineer/internal/usecase/getuser/model"

type dbUser struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	DateOfBirth string  `db:"date_of_birth"`
	Address     string  `db:"address"`
	Description *string `db:"description"`
	Latitude    float64 `db:"latitude"`
	Longitude   float64 `db:"longitude"`
}

func (db dbUser) toModel() model.User {
	return model.User{
		ID:          db.ID,
		Name:        db.Name,
		DateOfBirth: db.DateOfBirth,
		Address:     db.Address,
		Description: db.Description,
		Location: model.Point{
			Latitude:  db.Latitude,
			Longitude: db.Longitude,
		},
	}
}
