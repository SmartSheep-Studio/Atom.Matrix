package controllers

import (
	"code.smartsheep.studio/atom/matrix/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/middleware"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CloudSaveController struct {
	db         *gorm.DB
	gatekeeper *middleware.AuthMiddleware
}

func NewCloudSaveController(db *gorm.DB, gatekeeper *middleware.AuthMiddleware) *CloudSaveController {
	return &CloudSaveController{db, gatekeeper}
}

func (ctrl *CloudSaveController) Map(router *fiber.App) {
	router.Get(
		"/api/apps/:app/cloud-save",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:cloud-save"), hyperutils.GenPerms()),
		ctrl.get,
	)
	router.Put(
		"/api/apps/:app/cloud-save",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("update:cloud-save"), hyperutils.GenPerms()),
		ctrl.update,
	)
	router.Put(
		"/api/apps/:app/cloud-save/name",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("update:cloud-save.meta"), hyperutils.GenPerms()),
		ctrl.updateMeta,
	)
}

func (ctrl *CloudSaveController) get(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ?", c.Params("app")).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var library models.LibraryItem
	if err := ctrl.db.Where("app_id = ? AND user_id = ?", app.ID, u.ID).Preload("CloudSave").First(&library).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(fiber.StatusForbidden, "you haven't that app")
		} else {
			return hyperutils.ErrorParser(err)
		}
	}

	return c.JSON(library.CloudSave)
}

func (ctrl *CloudSaveController) update(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var req map[string]any
	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("you need provide a valid json format payload: %q", err))
	}

	var app models.App
	if err := ctrl.db.Where("slug = ?", c.Params("app")).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var library models.LibraryItem
	if err := ctrl.db.Where("app_id = ? AND user_id = ?", app.ID, u.ID).Preload("CloudSave").First(&library).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(fiber.StatusForbidden, "you haven't that app")
		} else {
			return hyperutils.ErrorParser(err)
		}
	}

	library.CloudSave.Payload = datatypes.JSON(data)

	if err := ctrl.db.Save(&library.CloudSave).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(library.CloudSave)
	}
}

func (ctrl *CloudSaveController) updateMeta(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var req struct {
		Name string `json:"name" validate:"required"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	var app models.App
	if err := ctrl.db.Where("slug = ?", c.Params("app")).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var library models.LibraryItem
	if err := ctrl.db.Where("app_id = ? AND user_id = ?", app.ID, u.ID).Preload("CloudSave").First(&library).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(fiber.StatusForbidden, "you haven't that app")
		} else {
			return hyperutils.ErrorParser(err)
		}
	}

	library.CloudSave.Name = req.Name

	if err := ctrl.db.Save(&library.CloudSave).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(library.CloudSave)
	}
}
