package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/internal/mocks"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var (
		testRadiusMeter      = float64(500)
		testUserID           = uuid.NewString()
		testGetNearbyFriends = model.Friends{
			model.Friend{
				ID:          "friend-1",
				Name:        "friend-1",
				DateOfBirth: "1999-12-12",
				Address:     "Orchid Road",
				Location: model.Point{
					Longitude: 123.44,
					Latitude:  23.44,
				},
			},
			model.Friend{
				ID:          "friend-2",
				Name:        "friend-2",
				DateOfBirth: "1999-12-12",
				Address:     "Klang Lama",
				Location: model.Point{
					Longitude: 100.44,
					Latitude:  22.44,
				},
			},
		}
	)
	testCases := []struct {
		name             string
		userID           string
		radiusMeter      float64
		expectedResult   *model.Friends
		expectedError    error
		repoGetNearbyErr error
		repoGetNearbyRes model.Friends
	}{
		{
			name:             "Success get nearby friends",
			userID:           testUserID,
			radiusMeter:      testRadiusMeter,
			expectedResult:   &testGetNearbyFriends,
			repoGetNearbyRes: testGetNearbyFriends,
		},
		{
			name:             "Repository get nearby friends error",
			userID:           testUserID,
			radiusMeter:      testRadiusMeter,
			repoGetNearbyErr: errors.New("repo error"),
			expectedError:    errors.New("repo error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)
			uc := UseCase{
				repo: mockRepo,
			}
			ctx := context.Background()
			mockRepo.EXPECT().
				Get(gomock.Any(), tc.userID, tc.radiusMeter).
				Times(1).
				Return(tc.repoGetNearbyRes, tc.repoGetNearbyErr)

			// Action
			actualResult, actualErr := uc.GetNearby(ctx, tc.userID, tc.radiusMeter)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, actualErr)
			}

			if tc.expectedResult != nil {
				assert.Equal(t, *tc.expectedResult, actualResult)
			}
		})
	}
}
