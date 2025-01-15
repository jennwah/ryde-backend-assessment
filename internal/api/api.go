package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jennwah/ryde-backend-engineer/docs"
	friendCtrl "github.com/jennwah/ryde-backend-engineer/internal/api/controller/friend"
	userCtrl "github.com/jennwah/ryde-backend-engineer/internal/api/controller/user"
	"github.com/jennwah/ryde-backend-engineer/internal/config"
)

type API struct {
	userController   userCtrl.Controller
	friendController friendCtrl.Controller
}

func New(
	db *sqlx.DB,
	logger *slog.Logger,
	cfg config.Config,
) *API {
	return &API{
		userController:   userCtrl.New(db, logger, cfg),
		friendController: friendCtrl.New(db, logger, cfg),
	}
}

func (a API) RegisterHandlers(r gin.IRouter) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	userV1ApiGroup := r.Group("api/v1/users")
	{
		userV1ApiGroup.POST("", a.userController.Create)
		userV1ApiGroup.GET("/:id", a.userController.GetByID)
		userV1ApiGroup.PATCH("/:id", a.userController.Patch)
		userV1ApiGroup.DELETE("/:id", a.userController.Delete)
	}

	friendV1ApiGroup := r.Group("api/v1/friends")
	{
		friendV1ApiGroup.POST("", a.friendController.Create)
		friendV1ApiGroup.GET("/nearby", a.friendController.GetNearby)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
