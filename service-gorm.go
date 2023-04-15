package service

import (
	"embed"

	"github.com/itbasis/go-service/db"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (receiver *Service) InitDB(dbEmbedMigrations *embed.FS) *Service {
	if receiver.config.DbGormDisabled {
		log.Info().Msg(gormIsDisabled)

		return receiver
	}

	newDB, err := db.NewDB(dbEmbedMigrations)
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	receiver.gorm = newDB

	return receiver
}

func (receiver *Service) GetGorm() *gorm.DB {
	return receiver.GetGormWithEmbeddedMigrations(nil)
}

func (receiver *Service) GetGormWithEmbeddedMigrations(dbEmbedMigrations *embed.FS) *gorm.DB {
	if receiver.gorm != nil {
		return receiver.gorm.GetGorm()
	}

	if receiver.config.DbGormDisabled {
		log.Error().Msg(gormIsDisabled)

		return nil
	}

	if receiver.gorm == nil {
		receiver.InitDB(dbEmbedMigrations)
	}

	return receiver.gorm.GetGorm()
}
