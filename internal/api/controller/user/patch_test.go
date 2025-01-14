package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/jennwah/ryde-backend-engineer/internal/api/controller/user/mocks"
	"github.com/jennwah/ryde-backend-engineer/internal/config"
)

func TestController_Patch(t *testing.T) {
	var (
		testUUID         = uuid.NewString()
		validLongitude   = 152.22
		validLatitude    = 45.22
		invalidLongitude = -270.22
		invalidLatitude  = -123.22
	)

	tests := []struct {
		name            string
		request         PatchRequest
		isUseCaseCalled bool
		patchErr        error
		expectedStatus  int
	}{
		{
			name: "Success patch user with valid request",
			request: PatchRequest{
				Name:        "Test User",
				DateOfBirth: "1997-12-17",
				Address:     "Test Address",
				Latitude:    validLatitude,
				Longitude:   validLongitude,
			},
			isUseCaseCalled: true,
			expectedStatus:  http.StatusNoContent,
		},
		{
			name: "Patch user request - invalid date of birth",
			request: PatchRequest{
				Name:        "Test User",
				DateOfBirth: "12-invalid-date",
				Address:     "Test Address",
				Latitude:    validLatitude,
				Longitude:   validLongitude,
			},
			isUseCaseCalled: false,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name: "Patch user request - invalid latitude",
			request: PatchRequest{
				Name:        "Test User",
				DateOfBirth: "12-invalid-date",
				Address:     "Test Address",
				Latitude:    invalidLatitude,
				Longitude:   validLongitude,
			},
			isUseCaseCalled: false,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name: "Patch user request - invalid longitude",
			request: PatchRequest{
				Name:        "Test User",
				DateOfBirth: "12-invalid-date",
				Address:     "Test Address",
				Latitude:    validLatitude,
				Longitude:   invalidLongitude,
			},
			isUseCaseCalled: false,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name: "Patch user request - internal server error",
			request: PatchRequest{
				Name:        "Test User",
				DateOfBirth: "1997-12-17",
				Address:     "Test Address",
				Latitude:    validLatitude,
				Longitude:   validLongitude,
			},
			isUseCaseCalled: true,
			patchErr:        errors.New("some internal error"),
			expectedStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPatchUseCase := mocks.NewMockpatchUser(ctrl)
			if tt.isUseCaseCalled {
				mockPatchUseCase.EXPECT().
					Patch(gomock.Any(), gomock.Any()).
					Return(tt.patchErr)
			}

			controller := Controller{
				uc: UseCases{
					PatchUserUseCase: mockPatchUseCase,
				},
				logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
				cfg:    config.Config{},
			}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.AddParam("id", testUUID)

			// Prepare request body
			body, err := json.Marshal(tt.request)
			assert.NoError(t, err)
			ctx.Request = httptest.NewRequest(http.MethodPatch, "/", bytes.NewBuffer(body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			controller.Patch(ctx)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
