package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Controller) Get(c *gin.Context) {
	c.String(http.StatusOK, "Welcome Gin Server")
}
