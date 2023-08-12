package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"code.smartsheep.studio/atom/matrix/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/matrix/pkg/server/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ExploreController struct {
	db      *gorm.DB
	conn    *subapps.HeLiCoPtErConnection
	service *services.ExploreService
}

func NewExploreController(db *gorm.DB, conn *subapps.HeLiCoPtErConnection, service *services.ExploreService) *ExploreController {
	return &ExploreController{db, conn, service}
}

func (ctrl *ExploreController) Map(router *fiber.App) {
	router.Get("/api/explore/apps", ctrl.apps)
	router.Get("/api/explore/apps/:app", ctrl.app)
	router.Get("/api/explore/apps/:app/posts", ctrl.posts)
	router.Get("/api/explore/apps/:app/releases", ctrl.releases)
}

func (ctrl *ExploreController) apps(c *fiber.Ctx) error {
	items, err := ctrl.service.ExploreApps()
	if err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *ExploreController) app(c *fiber.Ctx) error {
	var app models.App
	id, _ := c.ParamsInt("app", 0)
	if err := ctrl.db.Where("slug = ? OR id = ?", c.Params("app"), id).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	items, err := ctrl.service.ExploreApp(app.ID)
	if err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *ExploreController) posts(c *fiber.Ctx) error {

	var app models.App
	id, _ := c.ParamsInt("app", 0)
	if err := ctrl.db.Where("slug = ? OR id = ?", c.Params("app"), id).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	items, err := ctrl.service.ExplorePosts(app.ID)
	if err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(items)
	}
}

func (ctrl *ExploreController) releases(c *fiber.Ctx) error {

	var app models.App
	id, _ := c.ParamsInt("app", 0)
	if err := ctrl.db.Where("slug = ? OR id = ?", c.Params("app"), id).First(&app).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	items, err := ctrl.service.ExploreReleases(app.ID)
	if err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(items)
	}
}
