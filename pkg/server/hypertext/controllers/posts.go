package controllers

import (
	"code.smartsheep.studio/atom/matrix/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PostController struct {
	db         *gorm.DB
	gatekeeper *middleware.AuthMiddleware
}

func NewPostController(db *gorm.DB, gatekeeper *middleware.AuthMiddleware) *PostController {
	return &PostController{db, gatekeeper}
}

func (ctrl *PostController) Map(router *fiber.App) {
	router.Get(
		"/api/apps/:app/posts",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:apps.posts"), hyperutils.GenPerms("apps.posts.read")),
		ctrl.list,
	)
	router.Get(
		"/api/apps/:app/posts/:post",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:apps.posts"), hyperutils.GenPerms("apps.posts.read")),
		ctrl.get,
	)
	router.Post(
		"/api/apps/:app/posts",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("create:apps.posts"), hyperutils.GenPerms("apps.posts.create")),
		ctrl.create,
	)
	router.Put(
		"/api/apps/:app/posts/:post",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("update:apps.posts"), hyperutils.GenPerms("apps.posts.update")),
		ctrl.update,
	)
	router.Delete(
		"/api/apps/:app/posts/:post",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("delete:apps.posts"), hyperutils.GenPerms("apps.posts.delete")),
		ctrl.delete,
	)
}

func (ctrl *PostController) list(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var posts []models.Post
	if err := ctrl.db.Where("app_id = ?", app.ID).Order("created_at desc").Find(&posts).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(posts)
	}
}

func (ctrl *PostController) get(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var post models.Post
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("post"), app.ID).First(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(post)
	}
}

func (ctrl *PostController) create(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Title       string   `json:"title" validate:"required"`
		Type        string   `json:"type" validate:"required"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
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
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(post)
	}
}

func (ctrl *PostController) update(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var req struct {
		Title       string   `json:"title" validate:"required"`
		Type        string   `json:"type" validate:"required"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
		IsPublished bool     `json:"is_published"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	var post models.Post
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("post"), app.ID).First(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	post.Title = req.Title
	post.Type = req.Type
	post.Content = req.Content
	post.Tags = datatypes.NewJSONSlice(req.Tags)
	post.IsPublished = req.IsPublished

	if err := ctrl.db.Save(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(post)
	}
}

func (ctrl *PostController) delete(c *fiber.Ctx) error {

	u := c.Locals("matrix-id").(*models.Account)

	var app models.App
	if err := ctrl.db.Where("slug = ? AND account_id = ?", c.Params("app"), u.ID).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var post models.Post
	if err := ctrl.db.Where("slug = ? AND app_id = ?", c.Params("post"), app.ID).First(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	if err := ctrl.db.Delete(&post).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
