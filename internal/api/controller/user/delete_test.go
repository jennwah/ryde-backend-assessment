package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jennwah/ryde-backend-engineer/internal/api/controller/user/mocks"
	"github.com/jennwah/ryde-backend-engineer/internal/config"
	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
)

func TestController_Delete(t *testing.T) {
	var (
		testUUID = uuid.NewString()
	)

	tests := []struct {
		name            string
		pathParam       string
		isUseCaseCalled bool
		deleteErr       error
		expectedStatus  int
	}{
		{
			name:            "Success delete user by valid UUID",
			pathParam:       testUUID,
			isUseCaseCalled: true,
			expectedStatus:  http.StatusNoContent,
		},
		{
			name:           "Failed delete user with invalid id",
			pathParam:      "random-invalid-uuid-string",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "Not found user",
			pathParam:       testUUID,
			isUseCaseCalled: true,
			deleteErr:       postgresql.ErrNotFound,
			expectedStatus:  http.StatusNotFound,
		},
		{
			name:            "Internal server error",
			pathParam:       testUUID,
			isUseCaseCalled: true,
			deleteErr:       errors.New("internal error"),
			expectedStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDeleteUseCase := mocks.NewMockdeleteUser(ctrl)
			if tt.isUseCaseCalled {
				mockDeleteUseCase.EXPECT().
					Delete(gomock.Any(), tt.pathParam).
					Return(tt.deleteErr)
			}

			controller := Controller{
				uc: UseCases{
					DeleteUserUseCase: mockDeleteUseCase,
				},
				logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
				cfg:    config.Config{},
			}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.AddParam("id", tt.pathParam)

			controller.Delete(ctx)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
