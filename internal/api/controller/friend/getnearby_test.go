package friend

import (
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
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/model"
)

func TestController_GetNearby(t *testing.T) {
	var (
		testUUID             = uuid.NewString()
		testGetNearbyFriends = model.Friends{
			model.Friend{
				ID:          "friend-1",
				Name:        "friend-1",
				DateOfBirth: "1999-12-12",
				Address:     "Orchid Road",
				Location: model.Point{
					Longitude: 123.44,
					Latitude:  23.44,
				},
			},
			model.Friend{
				ID:          "friend-2",
				Name:        "friend-2",
				DateOfBirth: "1999-12-12",
				Address:     "Klang Lama",
				Location: model.Point{
					Longitude: 100.44,
					Latitude:  22.44,
				},
			},
		}
		testGetNearbyResp = GetNearbyResp{
			Friends: []Friend{
				{
					ID:          "friend-1",
					Name:        "friend-1",
					DateOfBirth: "1999-12-12",
					Address:     "Orchid Road",
					Longitude:   123.44,
					Latitude:    23.44,
				},
				{
					ID:          "friend-2",
					Name:        "friend-2",
					DateOfBirth: "1999-12-12",
					Address:     "Klang Lama",
					Longitude:   100.44,
					Latitude:    22.44,
				},
			},
		}
	)

	tests := []struct {
		name                  string
		queryParamRadiusMeter string
		userIDHeader          string
		isUseCaseCalled       bool
		getRes                model.Friends
		getErr                error
		expectedStatus        int
		expectedRes           *GetNearbyResp
	}{
		{
			name:                  "Success - get nearby friends",
			queryParamRadiusMeter: "5000.35",
			userIDHeader:          testUUID,
			isUseCaseCalled:       true,
			getRes:                testGetNearbyFriends,
			expectedStatus:        http.StatusOK,
			expectedRes:           &testGetNearbyResp,
		},
		{
			name:                  "Invalid header",
			queryParamRadiusMeter: "5000.35",
			userIDHeader:          "invalid-header",
			isUseCaseCalled:       false,
			expectedStatus:        http.StatusBadRequest,
		},
		{
			name:                  "Nearby friends - Not found",
			queryParamRadiusMeter: "5000.35",
			userIDHeader:          testUUID,
			isUseCaseCalled:       true,
			getRes:                model.Friends{},
			expectedStatus:        http.StatusNotFound,
		},
		{
			name:                  "Internal server error",
			queryParamRadiusMeter: "5000.35",
			userIDHeader:          testUUID,
			isUseCaseCalled:       true,
			getErr:                errors.New("internal error"),
			expectedStatus:        http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGetNearbyUseCase := mocks.NewMockgetNearbyFriends(ctrl)
			if tt.isUseCaseCalled {
				mockGetNearbyUseCase.EXPECT().
					GetNearby(gomock.Any(), tt.userIDHeader, gomock.Any()).
					Return(tt.getRes, tt.getErr)
			}

			controller := Controller{
				uc: UseCases{
					GetNearbyFriendsUseCase: mockGetNearbyUseCase,
				},
				logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
				cfg:    config.Config{},
			}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			ctx.Request.Header.Set("user_id", tt.userIDHeader)

			controller.GetNearby(ctx)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedRes != nil {
				var response GetNearbyResp
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, *tt.expectedRes, response)
			}
		})
	}
}
