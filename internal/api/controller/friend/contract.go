package friend

// MOCK: make mock_interface FILE_PATH=./internal/api/controller/friend/contract.go DESTINATION_FILE=./internal/api/controller/friend/mocks/mock.go PACKAGE=mocks
import (
	"context"

	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/model"
)

type createFriend interface {
	Create(ctx context.Context, userID, friendID string) error
}

type getNearbyFriends interface {
	GetNearby(ctx context.Context, userID string, radiusMeter float64) (model.Friends, error)
}
