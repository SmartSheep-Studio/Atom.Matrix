package controllers

import (
	"encoding/json"
	"errors"
	"fmt"

	"code.smartsheep.studio/atom/neutron/http/context"
	"github.com/gofiber/fiber/v2"

	"code.smartsheep.studio/atom/matrix/datasource/models"
	"code.smartsheep.studio/atom/matrix/http/middleware"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CloudSaveController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewCloudSaveController(db *gorm.DB, auth middleware.AuthHandler) *CloudSaveController {
	return &CloudSaveController{db, auth}
}

func (ctrl *CloudSaveController) Map(router *context.App) {
	router.Get("/api/apps/:app/cloud-save", ctrl.auth(true, "matrix.cloud-save.read"), ctrl.get)
	router.Put("/api/apps/:app/cloud-save", ctrl.auth(true, "matrix.cloud-save.update"), ctrl.update)
	router.Put("/api/apps/:app/cloud-save/name", ctrl.auth(true, "matrix.cloud-save.update.name"), ctrl.updateInfo)
}

func (ctrl *CloudSaveController) get(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ?", c.Params("app")).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var library models.LibraryItem
	if err := ctrl.db.Where("app_id = ? AND user_id = ?", app.ID, u.ID).Preload("CloudSave").First(&library).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(fiber.StatusForbidden, "you haven't that app")
		} else {
			return c.DbError(err)
		}
	}

	return c.JSON(library.CloudSave)
}

func (ctrl *CloudSaveController) update(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var req map[string]any
	if err := c.BindBody(&req); err != nil {
		return err
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("you need provide a valid json format payload: %q", err))
	}

	var app models.App
	if err := ctrl.db.Where("slug = ?", c.Params("app")).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var library models.LibraryItem
	if err := ctrl.db.Where("app_id = ? AND user_id = ?", app.ID, u.ID).Preload("CloudSave").First(&library).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(fiber.StatusForbidden, "you haven't that app")
		} else {
			return c.DbError(err)
		}
	}

	library.CloudSave.Payload = datatypes.JSON(data)

	if err := ctrl.db.Save(&library.CloudSave).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(library.CloudSave)
	}
}

func (ctrl *CloudSaveController) updateInfo(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var req struct {
		Name string `json:"name" validate:"required"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	var app models.App
	if err := ctrl.db.Where("slug = ?", c.Params("app")).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var library models.LibraryItem
	if err := ctrl.db.Where("app_id = ? AND user_id = ?", app.ID, u.ID).Preload("CloudSave").First(&library).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fiber.NewError(fiber.StatusForbidden, "you haven't that app")
		} else {
			return c.DbError(err)
		}
	}

	library.CloudSave.Name = req.Name

	if err := ctrl.db.Save(&library.CloudSave).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(library.CloudSave)
	}
}
