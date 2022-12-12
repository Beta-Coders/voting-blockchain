package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"votingblockchain/container"
	"votingblockchain/model"
)

type HealthController struct {
	container container.Container
}

func NewHealthController(container container.Container) *HealthController {
	return &HealthController{container: container}
}

// GetHealth - check server health
func (controller *HealthController) GetHealth(c echo.Context) error {
	msg := &model.Response{Message: "healthy"}
	return c.JSON(http.StatusOK, msg)
}
