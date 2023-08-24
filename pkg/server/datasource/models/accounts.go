package models

type Account struct {
	Model

	Nickname   string        `json:"nickname"`
	Experience int64         `json:"experience"`
	Library    []LibraryItem `json:"library" gorm:"foreignKey:AccountID"`
	UserID     uint          `json:"user_id"`
}

type LibraryItem struct {
	Model

	PlayTime  int64     `json:"play_time"`
	CloudSave CloudSave `json:"cloud_save" gorm:"foreignKey:LibraryID"`
	AccountID uint      `json:"account_id"`
	AppID     uint      `json:"app_id"`
}
