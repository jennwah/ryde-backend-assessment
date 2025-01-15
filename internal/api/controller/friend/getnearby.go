package friend

import (
	"github.com/gin-gonic/gin"
	"github.com/jennwah/ryde-backend-engineer/internal/api/controller"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/getnearbyfriends/model"
	"log/slog"
	"net/http"
	"strconv"
)

const DefaultRadiusMeter = 5000

// GetNearby GetNearbyFriends godoc
// @BasePath     /api/v1
// @Summary      Get nearby friends with radius meter specified
// @Description  Retrieve nearby friends from database, within radius meter specified
// @Tags         friends
// @Produce      json
// @Param        radius_meter  query    number  false  "Radius in meters (default: 5000)"  default(5000)
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  friend.GetNearbyResp
// @Failure      400  {object}  controller.ErrorResponseModel "Bad request"
// @Failure      404  {object}  controller.ErrorResponseModel "Not found"
// @Failure      500  {object}  controller.ErrorResponseModel "Internal server error"
// @Router       /friends/nearby [get]
func (h Controller) GetNearby(c *gin.Context) {
	var header Header

	if err := c.ShouldBindHeader(&header); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, controller.BadRequestErrorResp(err))
		return
	}

	radiusMeter := c.DefaultQuery("radius_meter", "0")
	parsedRadiusMeter, err := strconv.ParseFloat(radiusMeter, 64)
	if err != nil {
		parsedRadiusMeter = DefaultRadiusMeter
	}

	friends, err := h.uc.GetNearbyFriendsUseCase.GetNearby(c, header.UserId, parsedRadiusMeter)
	if err != nil {
		h.logger.Error("get user's nearby friends err", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, controller.InternalServerErrorResp())
		return
	}

	if len(friends) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, controller.NotFoundErrorResp())
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, getNearbyResp(friends))
}

func getNearbyResp(friends model.Friends) GetNearbyResp {
	friendsResp := make([]Friend, len(friends))

	for i, f := range friends {
		friendsResp[i] = Friend{
			ID:          f.ID,
			Name:        f.Name,
			Description: f.Description,
			Address:     f.Address,
			Latitude:    f.Location.Latitude,
			Longitude:   f.Location.Longitude,
			DateOfBirth: f.DateOfBirth,
		}
	}

	return GetNearbyResp{
		friendsResp,
	}
}
