package models

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	if err := db.AutoMigrate(
		&Account{},
		&App{},
		&LibraryItem{},
		&Post{},
		&Release{},
		&CloudSave{},
	); err != nil {
		log.Fatal().Err(err).Msg("Error when migrating database")
	}
}
