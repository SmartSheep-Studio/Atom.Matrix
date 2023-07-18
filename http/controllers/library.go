package controllers

import (
	"errors"

	"code.smartsheep.studio/atom/neutron/http/context"
	"github.com/gofiber/fiber/v2"

	"code.smartsheep.studio/atom/matrix/datasource/models"
	"code.smartsheep.studio/atom/matrix/http/middleware"
	"code.smartsheep.studio/atom/neutron/toolbox"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LibraryController struct {
	db   *gorm.DB
	conn *toolbox.ExternalServiceConnection
	auth middleware.AuthHandler
}

func NewLibraryController(db *gorm.DB, conn *toolbox.ExternalServiceConnection, auth middleware.AuthHandler) *LibraryController {
	return &LibraryController{db, conn, auth}
}

func (ctrl *LibraryController) Map(router *context.App) {
	router.Get("/api/library", ctrl.auth(true, "matrix.library.read"), ctrl.list)
	router.Get("/api/library/own", ctrl.auth(true, "matrix.library.read"), ctrl.doesOwn)
	router.Post("/api/library/add", ctrl.auth(true, "matrix.library.add"), ctrl.add)
}

func (ctrl *LibraryController) list(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var items []models.LibraryItem
	if err := ctrl.db.Where("account_id = ?", u.ID).Find(&items).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *LibraryController) doesOwn(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)
	target := c.Query("app")

	var app models.App
	if err := ctrl.db.Where("slug = ?", target).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var libraryCount int64
	if err := ctrl.db.Model(&models.LibraryItem{}).Where("account_id = ? AND app_id = ?", u.ID, app.ID).Count(&libraryCount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.DbError(err)
	} else if libraryCount <= 0 {
		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func (ctrl *LibraryController) add(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var req struct {
		App string `json:"app" validate:"required"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	var app models.App
	if err := ctrl.db.Where("slug = ?", req.App).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var libraryCount int64
	if err := ctrl.db.Model(&models.LibraryItem{}).Where("account_id = ? AND app_id = ?", u.ID, app.ID).Count(&libraryCount).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.DbError(err)
		}
	} else if libraryCount > 0 {
		return fiber.NewError(fiber.StatusForbidden, "already in the library")
	}

	item := models.LibraryItem{
		AccountID: u.ID,
		AppID:     app.ID,
		CloudSave: models.CloudSave{
			Name:    u.Nickname,
			Payload: datatypes.JSON([]byte("{}")),
		},
	}

	if err := ctrl.db.Save(&item).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(item)
	}
}
