package friend

import (
	"errors"
	"github.com/jennwah/ryde-backend-engineer/internal/pkg/postgresql"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jennwah/ryde-backend-engineer/internal/api/controller"
)

// Create CreateFriend godoc
// @BasePath     /api/v1
// @Summary      Create a new user friend pair
// @Description  Creates a new user friend pair record in the database.
// @Tags         friends
// @Accept       json
// @Param        friendHeader  header friend.Header  true "friend controller header"
// @Param        createRequest body friend.CreateReq true "Create user friend pair payload"
// @Success      201
// @Failure      400  {object}  controller.ErrorResponseModel "Bad request"
// @Failure      500  {object}  controller.ErrorResponseModel "Internal server error"
// @Router       /friends [post]
func (h Controller) Create(c *gin.Context) {
	var req CreateReq
	var header Header

	if err := c.ShouldBindHeader(&header); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, controller.BadRequestErrorResp(err))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, controller.BadRequestErrorResp(err))
		return
	}

	err := h.uc.CreateFriendUseCase.Create(c, header.UserId, req.FriendID)
	if errors.Is(err, postgresql.ErrFriendshipAlreadyExist) {
		c.AbortWithStatus(http.StatusCreated)
		return
	}

	if err != nil {
		h.logger.Error("create user friend pair err", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, controller.InternalServerErrorResp())
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}
