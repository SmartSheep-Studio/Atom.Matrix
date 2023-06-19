package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
	"repo.smartsheep.studio/atom/matrix/server/middleware"
	tmodels "repo.smartsheep.studio/atom/nucleus/datasource/models"
	"repo.smartsheep.studio/atom/nucleus/utils"
)

type AppController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewAppController(db *gorm.DB, auth middleware.AuthHandler) *AppController {
	return &AppController{db, auth}
}

func (ctrl *AppController) Map(router *fiber.App) {
	router.Get("/api/apps", ctrl.auth(true), ctrl.list)
	router.Get("/api/apps/:app", ctrl.auth(true), ctrl.get)
	router.Post("/api/apps", ctrl.auth(true), ctrl.create)
	router.Put("/api/apps/:app", ctrl.auth(true), ctrl.update)
	router.Delete("/api/apps/:app", ctrl.auth(true), ctrl.delete)
}

func (ctrl *AppController) list(c *fiber.Ctx) error {
	u := c.Locals("principal").(tmodels.User)

	var shops []models.MatrixApp
	if err := ctrl.db.Where("user_id = ?", u.ID).Find(&shops).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(shops)
	}
}

func (ctrl *AppController) get(c *fiber.Ctx) error {
	u := c.Locals("principal").(tmodels.User)

	var shop models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND user_id = ?", c.Params("app"), u.ID).First(&shop).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(shop)
	}
}

func (ctrl *AppController) create(c *fiber.Ctx) error {
	u := c.Locals("principal").(tmodels.User)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Details     string   `json:"details"`
		Url         string   `json:"url"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := utils.ParseRequestBody(c, &req); err != nil {
		return err
	}

	app := models.MatrixApp{
		Slug:        req.Slug,
		Url:         req.Url,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		Name:        req.Name,
		Description: req.Description,
		Details:     req.Details,
		IsPublished: req.IsPublished,
		UserID:      u.ID,
	}

	if err := ctrl.db.Save(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) update(c *fiber.Ctx) error {
	u := c.Locals("principal").(tmodels.User)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Details     string   `json:"details"`
		Url         string   `json:"url"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := utils.ParseRequestBody(c, &req); err != nil {
		return err
	}

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND user_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	app.Url = req.Url
	app.Tags = datatypes.NewJSONSlice(req.Tags)
	app.Slug = req.Slug
	app.Name = req.Name
	app.Description = req.Description
	app.Details = req.Details
	app.IsPublished = req.IsPublished

	if err := ctrl.db.Save(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.JSON(app)
	}
}

func (ctrl *AppController) delete(c *fiber.Ctx) error {
	u := c.Locals("principal").(tmodels.User)

	var app models.MatrixApp
	if err := ctrl.db.Where("slug = ? AND user_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	}

	if err := ctrl.db.Delete(&app).Error; err != nil {
		return utils.ParseDataSourceError(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
