package service

import (
	"embed"

	"github.com/itbasis/go-service/db"
	"github.com/rs/zerolog/log"
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

func (receiver *Service) GetGorm() *db.DB {
	if receiver.config.DbGormDisabled {
		log.Error().Msg(gormIsDisabled)

		return nil
	}

	return receiver.gorm
}
