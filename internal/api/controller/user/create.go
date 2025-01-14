package user

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jennwah/ryde-backend-engineer/internal/api/controller"
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/createuser/model"
)

// Create CreateUser godoc
// @BasePath     /api/v1
// @Summary      Create a new user
// @Description  Creates a new user record in the database.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        createRequest body user.CreateRequest true "Create user payload"
// @Success      201  {object}  user.CreateResp
// @Failure      400  {object}  controller.ErrorResponseModel "Bad request"
// @Failure      500  {object}  controller.ErrorResponseModel "Internal server error"
// @Router       /users [post]
func (h Controller) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, controller.BadRequestErrorResp(err))
		return
	}

	userID, err := h.uc.CreateUserUseCase.Create(c, model.User{
		Name:        req.Name,
		DateOfBirth: req.DateOfBirth,
		Address:     req.Address,
		Location: model.Point{
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		},
		Description: req.Description,
	})
	if err != nil {
		h.logger.Error("create user err", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, controller.InternalServerErrorResp())
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, CreateResp{
		ID: userID,
	})
}
