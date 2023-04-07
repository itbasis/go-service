package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v7"
	coreUtils "github.com/itbasis/go-core-utils"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Gorm *gorm.DB

	dbEmbedMigrations *embed.FS

	Config *Config

	connCredential       *Credential
	connSchemaCredential *Credential

	hostname string
}

func NewDB(dbEmbedMigrations *embed.FS) (*DB, error) {
	instance := &DB{
		dbEmbedMigrations:    dbEmbedMigrations,
		Config:               &Config{},
		connCredential:       &Credential{},
		connSchemaCredential: &Credential{},
	}

	if err := instance.readEnvironmentConfig(); err != nil {
		log.Error().Err(err).Msg("error reading data from environment")

		return nil, ErrCreateInstance
	}

	if err := instance.connectDB(); err != nil {
		log.Error().Err(err).Msg("database connection error")

		return nil, ErrConnectDB
	}

	if err := instance.migrationDB(); err != nil {
		log.Error().Err(err).Msg("database migration error")

		return instance, ErrDbMigration
	}

	return instance, nil
}

func (receiver *DB) readEnvironmentConfig() error {
	log.Info().Msg("reading environment...")

	if err := coreUtils.ReadEnvConfig(receiver.Config); err != nil {
		log.Error().Err(err).Msg("error reading configuration from environment")

		return ErrCreateInstance
	}

	if err := coreUtils.ReadEnvConfig(receiver.connCredential); err != nil {
		log.Error().Err(err).Msg("error reading credentials from environment")

		return ErrCreateInstance
	}

	// Schema Credential (for migrations)
	receiver.connSchemaCredential = &Credential{
		User:     receiver.connCredential.User,
		Password: receiver.connCredential.Password,
	}

	if err := coreUtils.ReadEnvConfig(receiver.connSchemaCredential, env.Options{Prefix: "SCHEMA_"}); err != nil {
		log.Error().Err(err).Msg("error reading credentials for schema from environment")

		return ErrCreateInstance
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Error().Err(err).Msg("failed to get hostname")

		return ErrCreateInstance
	}

	receiver.hostname = hostname

	log.Info().Msg("complete.")

	return nil
}

func (receiver *DB) getDSN(credential Credential) string {
	dbConn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s application_name=%s",
		receiver.Config.Host,
		receiver.Config.Port,
		credential.User,
		credential.Password,
		receiver.Config.Name,
		receiver.Config.SslMode,
		receiver.hostname,
	)

	log.Trace().Msgf("dbConn: %s", dbConn)

	return dbConn
}

func (receiver *DB) connectDB() error {
	log.Info().Msg("connecting to database...")

	gormLogLevel := logger.Warn
	if log.Debug().Enabled() {
		gormLogLevel = logger.Info
	}

	gormDB, err := gorm.Open(
		postgres.Open(receiver.getDSN(*receiver.connCredential)),
		&gorm.Config{
			Logger: logger.Default.LogMode(gormLogLevel),
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("")

		return ErrConnectDB
	}

	db, err := gormDB.DB()
	if err != nil {
		log.Error().Err(err).Msg("")

		return ErrConnectDB
	}

	db.SetMaxIdleConns(receiver.Config.MaxIdleConnections)
	db.SetMaxOpenConns(receiver.Config.MaxOpenConnections)
	db.SetConnMaxLifetime(time.Duration(receiver.Config.MaxLifetimeInMinutes) * time.Minute)

	var version string
	if res := gormDB.Raw("SELECT VERSION();").First(&version); res.Error != nil {
		log.Error().Err(res.Error).Msg("")

		return ErrConnectDB
	}

	log.Info().Msgf("Database version: %s", version)

	receiver.Gorm = gormDB

	return nil
}

func (receiver *DB) migrationDB() error {
	migrationDir, err := receiver.prepareMigrationDB()
	if err != nil {
		log.Error().Err(err).Msg("database migration error")

		return err
	} else if len(migrationDir) == 0 {
		log.Warn().Msg("Source for database migration not found - migration skipped")

		return nil
	}

	log.Info().Msgf("Directory with database migrations: %s", migrationDir)

	// We connect to the database with the rights to edit the database schema
	sqlDb, err := sql.Open(receiver.Config.Dialect, receiver.getDSN(*receiver.connSchemaCredential))
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database to run migrations")

		return ErrDbMigration
	}

	defer func() {
		if err := sqlDb.Close(); err != nil {
			log.Error().Err(err).Msg("failed to disconnect from database")
		}
	}()

	if err = goose.Up(sqlDb, migrationDir, goose.WithAllowMissing()); err != nil {
		log.Error().Err(err).Msg("")

		return ErrDbMigration
	}

	log.Info().Msg("Database migration completed")

	return nil
}

func (receiver *DB) prepareMigrationDB() (string, error) {
	if err := goose.SetDialect(receiver.Config.Dialect); err != nil {
		return "", err //nolint:wrapcheck
	}

	if receiver.dbEmbedMigrations != nil {
		log.Debug().Msg("Configure Goose using embedded FS for migrations")

		goose.SetBaseFS(receiver.dbEmbedMigrations)

		return "migrations", nil
	}

	migrationDir := receiver.Config.GooseMigrationDir
	if len(migrationDir) != 0 {
		// Checking the availability of the directory with database migrations
		_, err := os.Stat(migrationDir)
		if errors.Is(err, os.ErrNotExist) {
			log.Error().Err(err).Msgf("Directory with database migrations not found: %s", migrationDir)

			return "", err //nolint:wrapcheck
		}

		log.Info().Msgf("Directory with database migrations: %s", migrationDir)
	}

	return migrationDir, nil
}
