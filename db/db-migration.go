package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/juju/zaputil/zapctx"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

func (receiver *DB) migrationDB() error {
	logger := zapctx.Default.Sugar()

	migrationDir, err := receiver.prepareMigrationDB()
	if err != nil {
		e := errors.Wrap(err, "database migration error")
		logger.Error(e)

		return e
	} else if len(migrationDir) == 0 {
		logger.Warn("Source for database migration not found - migration skipped")

		return nil
	}

	logger.Infof("Directory with database migrations: %s", migrationDir)

	// We connect to the database with the rights to edit the database schema
	sqlDb, err := sql.Open(receiver.Config.Dialect, receiver.getDSN(*receiver.connSchemaCredential))
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to database to run migrations: %w", err))

		return ErrDbMigration
	}

	defer func() {
		if e := sqlDb.Close(); e != nil {
			logger.Error(fmt.Errorf("failed to disconnect from database: %w", e))
		}
	}()

	if e := goose.Up(sqlDb, migrationDir, goose.WithAllowMissing()); e != nil {
		logger.Error(e)

		return ErrDbMigration
	}

	logger.Info("Database migration completed")

	return nil
}

func (receiver *DB) prepareMigrationDB() (string, error) {
	if err := goose.SetDialect(receiver.Config.Dialect); err != nil {
		return "", err //nolint:wrapcheck
	}

	logger := zapctx.Default.Sugar()

	if receiver.dbEmbedMigrations != nil {
		logger.Debug("Configure Goose using embedded FS for migrations")

		goose.SetBaseFS(receiver.dbEmbedMigrations)

		return "migrations", nil
	}

	migrationDir := receiver.Config.GooseMigrationDir
	if len(migrationDir) != 0 {
		// Checking the availability of the directory with database migrations
		_, err := os.Stat(migrationDir)
		if errors.Is(err, os.ErrNotExist) {
			logger.Error(fmt.Errorf("directory '%s' with database migrations not found: %w", migrationDir, err))

			return "", err //nolint:wrapcheck
		}

		logger.Infof("Directory with database migrations: %s", migrationDir)
	}

	return migrationDir, nil
}
