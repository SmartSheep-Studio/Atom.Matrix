package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"github.com/gofiber/fiber/v2"
)

type StatusController struct {
	conn *subapps.HeLiCoPtErConnection
}

func NewStatusController(conn *subapps.HeLiCoPtErConnection) *StatusController {
	return &StatusController{conn}
}

func (ctrl *StatusController) Map(router *fiber.App) {
	router.Get("/api/info", ctrl.configure)
}

func (ctrl *StatusController) configure(c *fiber.Ctx) error {

	return c.JSON(ctrl.conn)
}
