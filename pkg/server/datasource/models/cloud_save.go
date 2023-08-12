package models

import "gorm.io/datatypes"

type CloudSave struct {
	Model

	Name      string         `json:"name"`
	Payload   datatypes.JSON `json:"payload"`
	LibraryID uint           `json:"library_id"`
}
