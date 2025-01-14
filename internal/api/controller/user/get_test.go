package user

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/jennwah/ryde-backend-engineer/internal/api/controller/user/mocks"
	"github.com/jennwah/ryde-backend-engineer/internal/config"
	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getuser/model"
)

func TestController_GetByID(t *testing.T) {
	var (
		testTime = time.Now()
		testUUID = uuid.NewString()
		testUser = model.User{
			Name:        "test",
			DateOfBirth: "17-12-1997",
			Address:     "test addresss",
			Location: model.Point{
				Latitude:  12.22,
				Longitude: 100.22,
			},
			CreatedAt: testTime,
		}
	)

	tests := []struct {
		name             string
		pathParam        string
		isUseCaseCalled  bool
		getResult        model.User
		getErr           error
		expectedStatus   int
		expectedResponse *GetResp
	}{
		{
			name:            "Success get user by valid UUID",
			pathParam:       testUUID,
			isUseCaseCalled: true,
			getResult:       testUser,
			getErr:          nil,
			expectedStatus:  http.StatusOK,
			expectedResponse: &GetResp{
				Name:        testUser.Name,
				DateOfBirth: testUser.DateOfBirth,
				Address:     testUser.Address,
				Latitude:    testUser.Location.Latitude,
				Longitude:   testUser.Location.Longitude,
			},
		},
		{
			name:           "Failed get user with invalid id",
			pathParam:      "random-invalid-uuid-string",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "Not found user",
			pathParam:       testUUID,
			isUseCaseCalled: true,
			getErr:          postgresql.ErrNotFound,
			expectedStatus:  http.StatusNotFound,
		},
		{
			name:            "Internal server error",
			pathParam:       testUUID,
			isUseCaseCalled: true,
			getErr:          errors.New("internal error"),
			expectedStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGetUseCase := mocks.NewMockgetUser(ctrl)
			if tt.isUseCaseCalled {
				mockGetUseCase.EXPECT().
					Get(gomock.Any(), tt.pathParam).
					Return(tt.getResult, tt.getErr)
			}

			controller := Controller{
				uc: UseCases{
					GetUserUseCase: mockGetUseCase,
				},
				logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
				cfg:    config.Config{},
			}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.AddParam("id", tt.pathParam)

			controller.GetByID(ctx)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedResponse != nil {
				var response GetResp
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, *tt.expectedResponse, response)
			}
		})
	}
}
