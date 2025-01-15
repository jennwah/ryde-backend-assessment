package repository

import "github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/model"

type dbFriend struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	DateOfBirth string  `db:"date_of_birth"`
	Address     string  `db:"address"`
	Description *string `db:"description"`
	Latitude    float64 `db:"latitude"`
	Longitude   float64 `db:"longitude"`
}

type dbFriends []dbFriend

func (db dbFriends) toModel() model.Friends {
	friends := make(model.Friends, len(db))

	for i, dbf := range db {
		friends[i] = model.Friend{
			ID:          dbf.ID,
			Address:     dbf.Address,
			Name:        dbf.Name,
			DateOfBirth: dbf.DateOfBirth,
			Description: dbf.Description,
			Location: model.Point{
				Latitude:  dbf.Latitude,
				Longitude: dbf.Longitude,
			},
		}
	}

	return friends
}
