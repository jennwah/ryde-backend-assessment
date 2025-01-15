package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	mock_contract "github.com/jennwah/ryde-backend-engineer/internal/usecase/createfriend/internal/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	var (
		testUserID = uuid.NewString()
	)
	testCases := []struct {
		name          string
		userID        string
		friendID      string
		expectedError error
		repoCreateErr error
	}{
		{
			name:     "Success create user friend pair",
			userID:   testUserID,
			friendID: testUserID,
		},
		{
			name:          "Repository create user friend error",
			userID:        testUserID,
			friendID:      testUserID,
			repoCreateErr: errors.New("repo error"),
			expectedError: errors.New("repo error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_contract.NewMockRepository(ctrl)
			uc := UseCase{
				repo: mockRepo,
			}
			ctx := context.Background()
			mockRepo.EXPECT().
				Create(gomock.Any(), tc.userID, tc.friendID).
				Times(1).
				Return(tc.repoCreateErr)

			// Action
			actualErr := uc.Create(ctx, tc.userID, tc.friendID)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, actualErr)
			} else {
				assert.Nil(t, actualErr)
			}
		})
	}
}
