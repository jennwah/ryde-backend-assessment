package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
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

func TestController_Create(t *testing.T) {
	var (
		testUUID         = uuid.NewString()
		validLongitude   = 152.22
		validLatitude    = 45.22
		invalidLongitude = -270.22
		invalidLatitude  = -123.22
	)

	tests := []struct {
		name             string
		request          CreateRequest
		isUseCaseCalled  bool
		createResult     string
		createErr        error
		expectedStatus   int
		expectedResponse *CreateResp
	}{
		{
			name: "Success create user with valid request",
			request: CreateRequest{
				Name:        "Test User",
				DateOfBirth: "1997-12-17",
				Address:     "Test Address",
				Latitude:    validLatitude,
				Longitude:   validLongitude,
			},
			isUseCaseCalled: true,
			createResult:    testUUID,
			createErr:       nil,
			expectedStatus:  http.StatusCreated,
			expectedResponse: &CreateResp{
				ID: testUUID,
			},
		},
		{
			name: "Create user request - invalid date of birth",
			request: CreateRequest{
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
			name: "Create user request - invalid latitude",
			request: CreateRequest{
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
			name: "Create user request - invalid longitude",
			request: CreateRequest{
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
			name: "Create user request - user already exists",
			request: CreateRequest{
				Name:        "Test User",
				DateOfBirth: "1997-12-17",
				Address:     "Test Address",
				Latitude:    validLatitude,
				Longitude:   validLongitude,
			},
			isUseCaseCalled: true,
			createResult:    "",
			createErr:       postgresql.ErrUserAlreadyExist,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name: "Create user request - internal server error",
			request: CreateRequest{
				Name:        "Test User",
				DateOfBirth: "1997-12-17",
				Address:     "Test Address",
				Latitude:    validLatitude,
				Longitude:   validLongitude,
			},
			isUseCaseCalled: true,
			createResult:    "",
			createErr:       errors.New("some internal error"),
			expectedStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCreateUseCase := mocks.NewMockcreateUser(ctrl)
			if tt.isUseCaseCalled {
				mockCreateUseCase.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(tt.createResult, tt.createErr)
			}

			controller := Controller{
				uc: UseCases{
					CreateUserUseCase: mockCreateUseCase,
				},
				logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
				cfg:    config.Config{},
			}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			// Prepare request body
			body, err := json.Marshal(tt.request)
			assert.NoError(t, err)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			controller.Create(ctx)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedResponse != nil {
				var response CreateResp
				err = json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, *tt.expectedResponse, response)
			}
		})
	}
}
