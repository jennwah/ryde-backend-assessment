package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/deleteuser/internal/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	var (
		testUserID = uuid.NewString()
	)
	testCases := []struct {
		name              string
		userID            string
		expectedError     error
		repoDeleteUserErr error
	}{
		{
			name:   "Success delete user",
			userID: testUserID,
		},
		{
			name:              "Repository delete user error",
			userID:            testUserID,
			repoDeleteUserErr: errors.New("repo error"),
			expectedError:     errors.New("repo error"),
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
				Delete(gomock.Any(), tc.userID).
				Times(1).
				Return(tc.repoDeleteUserErr)

			// Action
			actualErr := uc.Delete(ctx, tc.userID)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, actualErr)
			} else {
				assert.Nil(t, actualErr)
			}
		})
	}
}
