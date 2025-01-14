package user

// MOCK: make mock_interface FILE_PATH=./internal/api/controller/user/contract.go DESTINATION_FILE=./internal/api/controller/user/mocks/mock.go PACKAGE=mocks
import (
	"context"

	createusermodel "github.com/jennwah/ryde-backend-engineer/internal/usecase/createuser/model"
	getusermodel "github.com/jennwah/ryde-backend-engineer/internal/usecase/getuser/model"
	patchusermodel "github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/model"
)

type createUser interface {
	Create(ctx context.Context, user createusermodel.User) (string, error)
}

type getUser interface {
	Get(ctx context.Context, id string) (getusermodel.User, error)
}

type patchUser interface {
	Patch(ctx context.Context, user patchusermodel.PatchUser) error
}
