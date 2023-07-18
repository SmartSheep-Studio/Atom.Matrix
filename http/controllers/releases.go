package controllers

import (
	"code.smartsheep.studio/atom/matrix/datasource/models"
	"code.smartsheep.studio/atom/matrix/http/middleware"
	"code.smartsheep.studio/atom/neutron/http/context"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReleaseController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewReleaseController(db *gorm.DB, auth middleware.AuthHandler) *ReleaseController {
	return &ReleaseController{db, auth}
}

func (ctrl *ReleaseController) Map(router *context.App) {
	router.Get("/api/apps/:app/releases", ctrl.auth(true, "matrix.apps.releases.read", "matrix.apps.releases.read"), ctrl.list)
	router.Get("/api/apps/:app/releases/:release", ctrl.auth(true, "matrix.apps.releases.read", "matrix.apps.releases.read"), ctrl.get)
	router.Post("/api/apps/:app/releases", ctrl.auth(true, "matrix.apps.releases.create", "matrix.apps.releases.create"), ctrl.create)
	router.Put("/api/apps/:app/releases/:release", ctrl.auth(true, "matrix.apps.releases.update", "matrix.apps.releases.update"), ctrl.update)
	router.Delete("/api/apps/:app/releases/:release", ctrl.auth(true, "matrix.apps.releases.delete", "matrix.apps.releases.delete"), ctrl.delete)
}

func (ctrl *ReleaseController) list(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var releases []models.Release
	if err := ctrl.db.Where("app_id = ?", app.ID).Order("created_at desc").Preload("Post").Find(&releases).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(releases)
	}
}

func (ctrl *ReleaseController) get(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var release models.Release
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("release"), app.ID).Preload("Post").First(&release).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(release)
	}
}

func (ctrl *ReleaseController) create(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var req struct {
		Slug        string                `json:"slug" validate:"required"`
		Name        string                `json:"name" validate:"required"`
		Type        string                `json:"type" validate:"required"`
		Description string                `json:"description"`
		Details     string                `json:"details"`
		Tags        []string              `json:"tags"`
		Options     models.ReleaseOptions `json:"options" validate:"required"`
		IsPublished bool                  `json:"is_published"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	release := models.Release{
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		Post: models.Post{
			Slug:        req.Slug,
			Type:        req.Type,
			Title:       req.Name,
			Content:     req.Details,
			Tags:        datatypes.NewJSONSlice(req.Tags),
			IsPublished: req.IsPublished,
			AppID:       app.ID,
		},
		Options:     datatypes.NewJSONType(req.Options),
		IsPublished: req.IsPublished,
		AppID:       app.ID,
	}

	if err := ctrl.db.Save(&release).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(release)
	}
}

func (ctrl *ReleaseController) update(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var req struct {
		Slug        string                `json:"slug" validate:"required"`
		Name        string                `json:"name" validate:"required"`
		Type        string                `json:"type" validate:"required"`
		Description string                `json:"description"`
		Details     string                `json:"details"`
		Tags        []string              `json:"tags"`
		Options     models.ReleaseOptions `json:"options" validate:"required"`
		IsPublished bool                  `json:"is_published"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	tx := ctrl.db.Begin()

	var release models.Release
	if err := tx.Where("slug = ? AND app_id = ?", c.Params("release"), app.ID).Preload("Post").First(&release).Error; err != nil {
		tx.Rollback()
		return c.DbError(err)
	} else {
		ctrl.db.Unscoped().Delete(&release.Post)
	}

	release.Slug = req.Slug
	release.Name = req.Name
	release.Description = req.Description
	release.IsPublished = req.IsPublished
	release.Options = datatypes.NewJSONType(req.Options)
	release.Post = models.Post{
		Slug:        req.Slug,
		Type:        req.Type,
		Title:       req.Name,
		Content:     req.Details,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		IsPublished: req.IsPublished,
		AppID:       app.ID,
	}

	if err := tx.Save(&release).Error; err != nil {
		tx.Rollback()
		return c.DbError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(release)
	}
}

func (ctrl *ReleaseController) delete(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var release models.Release
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("release"), app.ID).Preload("Post").First(&release).Error; err != nil {
		return c.DbError(err)
	}

	if err := ctrl.db.Delete(&release.Post).Error; err != nil {
		return c.DbError(err)
	}

	if err := ctrl.db.Delete(&release).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
