package models

import "gorm.io/datatypes"

const (
	PostMinorUpdateType  = "minor-update"
	PostHotfixUpdateType = "hotfix-update"
	PostMajorUpdateType  = "major-update"
	PostAnnouncementType = "announcement"
	PostBlogType         = "blog"
)

type Post struct {
	Model

	Slug        string                      `json:"slug" gorm:"uniqueIndex"`
	Type        string                      `json:"type"`
	Title       string                      `json:"title"`
	Cover       string                      `json:"cover"`
	Content     string                      `json:"content"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	ReleaseID   *uint                       `json:"release_id"`
	IsPublished bool                        `json:"is_published"`
	AppID       uint                        `json:"app_id"`
}
