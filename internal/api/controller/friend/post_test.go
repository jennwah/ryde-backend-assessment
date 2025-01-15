package friend

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

	"github.com/jennwah/ryde-backend-engineer/internal/api/controller/friend/mocks"
	"github.com/jennwah/ryde-backend-engineer/internal/config"
)

func TestController_Create(t *testing.T) {
	var (
		testUUID = uuid.NewString()
	)

	tests := []struct {
		name            string
		request         CreateReq
		userIDHeader    string
		isUseCaseCalled bool
		createErr       error
		expectedStatus  int
		expectedErr     error
	}{
		{
			name: "Success create user friend pair",
			request: CreateReq{
				FriendID: testUUID,
			},
			userIDHeader:    testUUID,
			isUseCaseCalled: true,
			createErr:       nil,
			expectedStatus:  http.StatusCreated,
		},
		{
			name: "Invalid header",
			request: CreateReq{
				FriendID: testUUID,
			},
			userIDHeader:    "invalid-header",
			isUseCaseCalled: false,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name: "Invalid friend ID",
			request: CreateReq{
				FriendID: "invalid-friend-id",
			},
			userIDHeader:    testUUID,
			isUseCaseCalled: false,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name: "Internal server error",
			request: CreateReq{
				FriendID: testUUID,
			},
			userIDHeader:    testUUID,
			isUseCaseCalled: true,
			createErr:       errors.New("some internal error"),
			expectedStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCreateUseCase := mocks.NewMockcreateFriend(ctrl)
			if tt.isUseCaseCalled {
				mockCreateUseCase.EXPECT().
					Create(gomock.Any(), tt.userIDHeader, tt.request.FriendID).
					Return(tt.createErr)
			}

			controller := Controller{
				uc: UseCases{
					CreateFriendUseCase: mockCreateUseCase,
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
			ctx.Request.Header.Set("user_id", tt.userIDHeader)

			controller.Create(ctx)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
