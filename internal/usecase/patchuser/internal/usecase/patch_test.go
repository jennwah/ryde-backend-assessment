package usecase

import (
	"context"
	"errors"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/internal/mocks"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPatch(t *testing.T) {
	var (
		testUser = model.PatchUser{
			Name: "test",
		}
	)
	testCases := []struct {
		name             string
		userPayload      model.PatchUser
		expectedError    error
		repoPatchUserErr error
	}{
		{
			name:        "Success patch user",
			userPayload: testUser,
		},
		{
			name:             "Repository patch user error",
			userPayload:      testUser,
			repoPatchUserErr: errors.New("repo error"),
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
				Patch(gomock.Any(), tc.userPayload).
				Times(1).
				Return(tc.repoPatchUserErr)

			// Action
			actualErr := uc.Patch(ctx, tc.userPayload)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, actualErr)
			} else {
				assert.Nil(t, tc.expectedError)
			}
		})
	}
}
