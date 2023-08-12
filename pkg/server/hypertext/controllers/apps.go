package controllers

import (
	"code.smartsheep.studio/atom/matrix/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AppController struct {
	db         *gorm.DB
	gatekeeper *middleware.AuthMiddleware
}

func NewAppController(db *gorm.DB, gatekeeper *middleware.AuthMiddleware) *AppController {
	return &AppController{db, gatekeeper}
}

func (ctrl *AppController) Map(router *fiber.App) {
	router.Get(
		"/api/apps",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:apps"), hyperutils.GenPerms("apps.read")),
		ctrl.list,
	)
	router.Get(
		"/api/apps/:app",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:apps"), hyperutils.GenPerms("apps.read")),
		ctrl.get,
	)
	router.Post(
		"/api/apps",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("create:apps"), hyperutils.GenPerms("apps.create")),
		ctrl.create,
	)
	router.Put(
		"/api/apps/:app",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("update:apps"), hyperutils.GenPerms("apps.update")),
		ctrl.update,
	)
	router.Delete(
		"/api/apps/:app",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("delete:apps"), hyperutils.GenPerms("apps.delete")),
		ctrl.delete,
	)
}

func (ctrl *AppController) list(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var apps []models.App
	if err := ctrl.db.Where("account_id = ?", u.ID).Find(&apps).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(apps)
	}
}

func (ctrl *AppController) get(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) create(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Details     string   `json:"details"`
		Url         string   `json:"url"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	app := models.App{
		Slug:        req.Slug,
		Url:         req.Url,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		Name:        req.Name,
		Description: req.Description,
		Details:     req.Details,
		IsPublished: req.IsPublished,
		AccountID:   u.ID,
	}

	if err := ctrl.db.Save(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) update(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Details     string   `json:"details"`
		Url         string   `json:"url"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	app.Url = req.Url
	app.Tags = datatypes.NewJSONSlice(req.Tags)
	app.Slug = req.Slug
	app.Name = req.Name
	app.Description = req.Description
	app.Details = req.Details
	app.IsPublished = req.IsPublished

	if err := ctrl.db.Save(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) delete(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	if err := ctrl.db.Delete(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
