package services

import (
	"gorm.io/gorm"
	"repo.smartsheep.studio/atom/matrix/datasource/models"
)

type ExploreService struct {
	db *gorm.DB
}

func NewExploreService(db *gorm.DB) *ExploreService {
	return &ExploreService{db}
}

func (v *ExploreService) ExploreApps() ([]models.MatrixApp, error) {
	tx := v.db.Where("is_published = true")

	var apps []models.MatrixApp
	if err := tx.Limit(100).Find(&apps).Error; err != nil {
		return nil, err
	} else {
		return apps, nil
	}
}

func (v *ExploreService) ExplorePosts(app uint) ([]models.MatrixPost, error) {
	tx := v.db.Where("is_published = true AND app_id = ?", app)

	var posts []models.MatrixPost
	if err := tx.Limit(100).Find(&posts).Error; err != nil {
		return nil, err
	} else {
		return posts, nil
	}
}

func (v *ExploreService) ExploreReleases(app uint) ([]models.MatrixRelease, error) {
	tx := v.db.Where("is_published = true AND app_id = ?", app)

	var releases []models.MatrixRelease
	if err := tx.Limit(100).Find(&releases).Error; err != nil {
		return nil, err
	} else {
		return releases, nil
	}
}
