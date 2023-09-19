package service

import (
	"context"
	"embed"

	"github.com/itbasis/go-service/v2/db"
	"github.com/juju/zaputil/zapctx"
	"gorm.io/gorm"
)

func (receiver *Service) InitDB(ctx context.Context, dbEmbedMigrations *embed.FS) *Service {
	log := zapctx.Logger(ctx).Sugar()

	if receiver.config.DbGormDisabled {
		log.Info(gormIsDisabled)

		return receiver
	}

	newDB, err := db.NewDB(ctx, dbEmbedMigrations)
	if err != nil {
		log.Panic(err)
	}

	receiver.gorm = newDB

	return receiver
}

func (receiver *Service) GetGorm(ctx context.Context) *gorm.DB {
	return receiver.GetGormWithEmbeddedMigrations(ctx, nil)
}

func (receiver *Service) GetGormWithEmbeddedMigrations(ctx context.Context, dbEmbedMigrations *embed.FS) *gorm.DB {
	if receiver.gorm != nil {
		return receiver.gorm.GetGorm()
	}

	if receiver.config.DbGormDisabled {
		zapctx.Logger(ctx).Info(gormIsDisabled)

		return nil
	}

	if receiver.gorm == nil {
		receiver.InitDB(ctx, dbEmbedMigrations)
	}

	return receiver.gorm.GetGorm()
}
