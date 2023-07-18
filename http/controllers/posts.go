package controllers

import (
	"code.smartsheep.studio/atom/matrix/datasource/models"
	"code.smartsheep.studio/atom/matrix/http/middleware"
	"code.smartsheep.studio/atom/neutron/http/context"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PostController struct {
	db   *gorm.DB
	auth middleware.AuthHandler
}

func NewPostController(db *gorm.DB, auth middleware.AuthHandler) *PostController {
	return &PostController{db, auth}
}

func (ctrl *PostController) Map(router *context.App) {
	router.Get("/api/apps/:app/posts", ctrl.auth(true, "matrix.apps.posts.read", "matrix.apps.posts.read"), ctrl.list)
	router.Get("/api/apps/:app/posts/:post", ctrl.auth(true, "matrix.apps.posts.read", "matrix.apps.posts.read"), ctrl.get)
	router.Post("/api/apps/:app/posts", ctrl.auth(true, "matrix.apps.posts.create", "matrix.apps.posts.create"), ctrl.create)
	router.Put("/api/apps/:app/posts/:post", ctrl.auth(true, "matrix.apps.posts.update", "matrix.apps.posts.update"), ctrl.update)
	router.Delete("/api/apps/:app/posts/:post", ctrl.auth(true, "matrix.apps.posts.delete", "matrix.apps.posts.delete"), ctrl.delete)
}

func (ctrl *PostController) list(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var posts []models.Post
	if err := ctrl.db.Where("app_id = ?", app.ID).Order("created_at desc").Find(&posts).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(posts)
	}
}

func (ctrl *PostController) get(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var post models.Post
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("post"), app.ID).First(&post).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(post)
	}
}

func (ctrl *PostController) create(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Title       string   `json:"title" validate:"required"`
		Type        string   `json:"type" validate:"required"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	post := models.Post{
		Slug:        req.Slug,
		Type:        req.Type,
		Title:       req.Title,
		Content:     req.Content,
		Tags:        datatypes.NewJSONSlice(req.Tags),
		AppID:       app.ID,
		IsPublished: req.IsPublished,
	}

	if err := ctrl.db.Save(&post).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(post)
	}
}

func (ctrl *PostController) update(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var req struct {
		Title       string   `json:"title" validate:"required"`
		Type        string   `json:"type" validate:"required"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := c.BindBody(&req); err != nil {
		return err
	}

	var post models.Post
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("post"), app.ID).First(&post).Error; err != nil {
		return c.DbError(err)
	}

	post.Title = req.Title
	post.Type = req.Type
	post.Content = req.Content
	post.Tags = datatypes.NewJSONSlice(req.Tags)
	post.IsPublished = req.IsPublished

	if err := ctrl.db.Save(&post).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.JSON(post)
	}
}

func (ctrl *PostController) delete(ctx *fiber.Ctx) error {
	c := &context.Ctx{Ctx: ctx}
	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return c.DbError(err)
	}

	var post models.Post
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("post"), app.ID).First(&post).Error; err != nil {
		return c.DbError(err)
	}

	if err := ctrl.db.Delete(&post).Error; err != nil {
		return c.DbError(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
