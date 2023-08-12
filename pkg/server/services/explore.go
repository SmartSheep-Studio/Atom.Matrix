package services

import (
	models "code.smartsheep.studio/atom/matrix/pkg/server/datasource/models"
	"gorm.io/gorm"
)

type ExploreService struct {
	db *gorm.DB
}

func NewExploreService(db *gorm.DB) *ExploreService {
	return &ExploreService{db}
}

func (v *ExploreService) ExploreApps() ([]models.App, error) {
	tx := v.db.Where("is_published = true")

	var apps []models.App
	if err := tx.Limit(100).Order("created_at desc").Find(&apps).Error; err != nil {
		return nil, err
	} else {
		return apps, nil
	}
}

func (v *ExploreService) ExploreApp(app uint) (models.App, error) {
	tx := v.db.Where("is_published = true AND id = ?", app)

	var details models.App
	if err := tx.Preload("Posts").Preload("Releases").First(&details).Error; err != nil {
		return details, err
	} else {
		return details, nil
	}
}

func (v *ExploreService) ExplorePosts(app uint) ([]models.Post, error) {
	tx := v.db.Where("is_published = true AND app_id = ?", app)

	var posts []models.Post
	if err := tx.Limit(100).Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	} else {
		return posts, nil
	}
}

func (v *ExploreService) ExploreReleases(app uint) ([]models.Release, error) {
	tx := v.db.Where("is_published = true AND app_id = ?", app)

	var releases []models.Release
	if err := tx.Limit(100).Order("created_at desc").Find(&releases).Error; err != nil {
		return nil, err
	} else {
		return releases, nil
	}
}
