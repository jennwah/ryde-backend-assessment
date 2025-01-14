package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/createuser/internal/mocks"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/createuser/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	var (
		testUser = model.User{
			Name: "test",
		}
		testUserID = uuid.NewString()
	)
	testCases := []struct {
		name              string
		userPayload       model.User
		expectedResult    string
		expectedError     error
		repoCreateUserErr error
		repoCreateUserRes string
	}{
		{
			name:              "Success create user",
			userPayload:       testUser,
			expectedResult:    testUserID,
			repoCreateUserRes: testUserID,
		},
		{
			name:              "Repository create user error",
			userPayload:       testUser,
			repoCreateUserRes: "",
			repoCreateUserErr: errors.New("repo error"),
			expectedError:     errors.New("repo error"),
			expectedResult:    "",
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
				Create(gomock.Any(), tc.userPayload).
				Times(1).
				Return(tc.repoCreateUserRes, tc.repoCreateUserErr)

			// Action
			actualResult, actualErr := uc.Create(ctx, tc.userPayload)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, actualErr)
			}
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}
