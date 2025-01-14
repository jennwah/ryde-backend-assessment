package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jennwah/ryde-backend-engineer/internal/api/controller"
	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
	"log/slog"
	"net/http"
)

// Delete DeleteUser godoc
// @BasePath     /api/v1
// @Summary      Delete an existing user
// @Description  Delete an existing user record from the database.
// @Tags         users
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      204
// @Failure      400  {object}  controller.ErrorResponseModel "Bad request"
// @Failure      404  {object}  controller.ErrorResponseModel "Not found"
// @Failure      500  {object}  controller.ErrorResponseModel "Internal server error"
// @Router       /users/{id} [delete]
func (h Controller) Delete(c *gin.Context) {
	userID := c.Param("id")
	_, err := uuid.Parse(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, controller.BadRequestErrorResp(err))
		return
	}

	err = h.uc.DeleteUserUseCase.Delete(c, userID)
	if errors.Is(err, postgresql.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, controller.NotFoundErrorResp())
		return
	}
	if err != nil {
		h.logger.Error("delete user err", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, controller.InternalServerErrorResp())
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
