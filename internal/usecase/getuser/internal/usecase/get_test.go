package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getuser/internal/mocks"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getuser/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var (
		testUser = model.User{
			Name: "test",
		}
		testUserID = uuid.NewString()
	)
	testCases := []struct {
		name           string
		userID         string
		expectedResult model.User
		expectedError  error
		repoGetUserErr error
		repoGetUserRes model.User
	}{
		{
			name:           "Success get user",
			userID:         testUserID,
			expectedResult: testUser,
			repoGetUserRes: testUser,
		},
		{
			name:           "Repository get user error",
			userID:         testUserID,
			repoGetUserErr: errors.New("repo error"),
			expectedError:  errors.New("repo error"),
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
				Get(gomock.Any(), tc.userID).
				Times(1).
				Return(tc.repoGetUserRes, tc.repoGetUserErr)

			// Action
			actualResult, actualErr := uc.Get(ctx, tc.userID)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, actualErr)
			}
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}
