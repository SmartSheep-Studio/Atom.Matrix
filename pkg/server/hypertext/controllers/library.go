package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"code.smartsheep.studio/atom/matrix/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/middleware"
	"errors"

	"github.com/gofiber/fiber/v2"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LibraryController struct {
	db         *gorm.DB
	conn       *subapps.HeLiCoPtErConnection
	gatekeeper *middleware.AuthMiddleware
}

func NewLibraryController(db *gorm.DB, conn *subapps.HeLiCoPtErConnection, gatekeeper *middleware.AuthMiddleware) *LibraryController {
	return &LibraryController{db, conn, gatekeeper}
}

func (ctrl *LibraryController) Map(router *fiber.App) {
	router.Get(
		"/api/library",
		ctrl.gatekeeper.Fn(true, hyperutils.GenPerms("read:library"), hyperutils.GenPerms()),
		ctrl.list,
	)
	router.Get(
		"/api/library/own",
		ctrl.gatekeeper.Fn(true, hyperutils.GenPerms("read:library"), hyperutils.GenPerms()),
		ctrl.doesOwn,
	)
	router.Post(
		"/api/library/add",
		ctrl.gatekeeper.Fn(true, hyperutils.GenPerms("update:library"), hyperutils.GenPerms()),
		ctrl.add,
	)
}

func (ctrl *LibraryController) list(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var items []models.LibraryItem
	if err := ctrl.db.Where("account_id = ?", u.ID).Find(&items).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *LibraryController) doesOwn(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)
	target := c.Query("app")

	var app models.App
	if err := ctrl.db.Where("slug = ?", target).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var libraryCount int64
	if err := ctrl.db.Model(&models.LibraryItem{}).Where("account_id = ? AND app_id = ?", u.ID, app.ID).Count(&libraryCount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return hyperutils.ErrorParser(err)
	} else if libraryCount <= 0 {
		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func (ctrl *LibraryController) add(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var req struct {
		App string `json:"app" validate:"required"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	var app models.App
	if err := ctrl.db.Where("slug = ?", req.App).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var libraryCount int64
	if err := ctrl.db.Model(&models.LibraryItem{}).Where("account_id = ? AND app_id = ?", u.ID, app.ID).Count(&libraryCount).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return hyperutils.ErrorParser(err)
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
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(item)
	}
}
