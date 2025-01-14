package user

import (
	"context"

	userdomain "github.com/jennwah/ryde-backend-engineer/internal/usecase/createuser/model"
)

type createUser interface {
	Create(ctx context.Context, user userdomain.User) (string, error)
}
