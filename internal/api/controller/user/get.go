package user

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jennwah/ryde-backend-engineer/internal/api/controller"
	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
)

// GetByID GetUserByID godoc
// @BasePath     /api/v1
// @Summary      Get user by ID
// @Description  Retrieve an existing user by UUID from database.
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  user.GetResp
// @Failure      400  {object}  controller.ErrorResponseModel "Bad request"
// @Failure      404  {object}  controller.ErrorResponseModel "Not found"
// @Failure      500  {object}  controller.ErrorResponseModel "Internal server error"
// @Router       /users/{id} [get]
func (h Controller) GetByID(c *gin.Context) {
	userID := c.Param("id")
	_, err := uuid.Parse(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, controller.BadRequestErrorResp(err))
		return
	}

	user, err := h.uc.GetUserUseCase.Get(c, userID)
	if errors.Is(err, postgresql.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, controller.NotFoundErrorResp())
		return
	}
	if err != nil {
		h.logger.Error("get user err", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, controller.InternalServerErrorResp())
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, GetResp{
		Name:        user.Name,
		DateOfBirth: user.DateOfBirth,
		Address:     user.Address,
		Description: user.Description,
		Latitude:    user.Location.Latitude,
		Longitude:   user.Location.Longitude,
	})
}
