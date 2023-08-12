package models

import "gorm.io/datatypes"

type App struct {
	Model

	Slug        string                      `json:"slug"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Url         string                      `json:"url"`
	Details     string                      `json:"details"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	IsPublished bool                        `json:"is_published"`
	Posts       []Post                      `json:"posts" gorm:"foreignKey:AppID"`
	Releases    []Release                   `json:"releases" gorm:"foreignKey:AppID"`
	Duplicates  []LibraryItem               `json:"duplicates" gorm:"foreignKey:AppID"`
	AccountID   uint                        `json:"account_id"`
}
