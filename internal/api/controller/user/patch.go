package user

import (
	"errors"
	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/model"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jennwah/ryde-backend-engineer/internal/api/controller"
)

// Patch PatchUser godoc
// @BasePath     /api/v1
// @Summary      Update an existing user
// @Description  Update an existing user record in the database.
// @Tags         users
// @Accept       json
// @Param        patchRequest body user.PatchRequest true "Update user payload"
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200
// @Failure      400  {object}  controller.ErrorResponseModel "Bad request"
// @Failure      404  {object}  controller.ErrorResponseModel "Not found"
// @Failure      500  {object}  controller.ErrorResponseModel "Internal server error"
// @Router       /users/{id} [patch]
func (h Controller) Patch(c *gin.Context) {
	userID := c.Param("id")
	_, err := uuid.Parse(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, controller.BadRequestErrorResp(err))
		return
	}

	var req PatchRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, controller.BadRequestErrorResp(err))
		return
	}

	err = h.uc.PatchUserUseCase.Patch(c, model.PatchUser{
		ID:          userID,
		Name:        req.Name,
		DateOfBirth: req.DateOfBirth,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Description: req.Description,
	})
	if errors.Is(err, postgresql.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, controller.NotFoundErrorResp())
		return
	}
	if err != nil {
		h.logger.Error("patch user err", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, controller.InternalServerErrorResp())
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
