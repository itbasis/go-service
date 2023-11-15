package db

import (
	"context"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
	itbasisCoreUtilsEnvReader "github.com/itbasis/go-core-utils/v2/env-reader"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type DB struct {
	gorm *gorm.DB

	dbEmbedMigrations *embed.FS

	Config *Config

	connCredential       *Credential
	connSchemaCredential *Credential

	hostname string
}

func NewDB(ctx context.Context, dbEmbedMigrations *embed.FS) (*DB, error) {
	instance := &DB{
		dbEmbedMigrations:    dbEmbedMigrations,
		Config:               &Config{},
		connCredential:       &Credential{},
		connSchemaCredential: &Credential{},
	}

	logger := zapctx.Logger(ctx).Sugar()

	if err := instance.readEnvironmentConfig(ctx); err != nil {
		logger.Error(fmt.Errorf("error reading data from environment: %w", err))

		return nil, ErrCreateInstance
	}

	if err := instance.connectDB(); err != nil {
		logger.Error(fmt.Errorf("database connection error: %w", err))

		return nil, ErrConnectDB
	}

	if err := instance.migrationDB(); err != nil {
		logger.Error(fmt.Errorf("database migration error: %w", err))

		return instance, ErrDbMigration
	}

	return instance, nil
}

func (receiver *DB) GetGorm() *gorm.DB { return receiver.gorm }

func (receiver *DB) readEnvironmentConfig(ctx context.Context) error {
	logger := zapctx.Default.Sugar()
	logger.Info("reading environment...")

	if err := itbasisCoreUtilsEnvReader.ReadEnvConfig(ctx, receiver.Config, nil); err != nil {
		logger.Error(fmt.Errorf("error reading configuration from environment: %w", err))

		return ErrCreateInstance
	}

	if err := itbasisCoreUtilsEnvReader.ReadEnvConfig(ctx, receiver.connCredential, nil); err != nil {
		logger.Error(fmt.Errorf("error reading credentials from environment: %w", err))

		return ErrCreateInstance
	}

	// Schema Credential (for migrations)
	receiver.connSchemaCredential = &Credential{
		User:     receiver.connCredential.User,
		Password: receiver.connCredential.Password,
	}

	if err := itbasisCoreUtilsEnvReader.ReadEnvConfig(ctx, receiver.connSchemaCredential, &env.Options{Prefix: "SCHEMA_"}); err != nil {
		logger.Error(fmt.Errorf("error reading credentials for schema from environment: %w", err))

		return ErrCreateInstance
	}

	hostname, err := os.Hostname()
	if err != nil {
		logger.Error(fmt.Errorf("failed to get hostname: %w", err))

		return ErrCreateInstance
	}

	receiver.hostname = hostname

	logger.Info("complete.")

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

	zapctx.Default.Sugar().Debugf("dbConn: %s", dbConn)

	return dbConn
}

func (receiver *DB) connectDB() error {
	logger := zapctx.Default.Sugar()
	logger.Info("connecting to database...")

	gormLogLevel := gormLogger.Warn
	if logger.Desugar().Core().Enabled(zapcore.DebugLevel) {
		gormLogLevel = gormLogger.Info
	}

	gormDB, err := gorm.Open(
		postgres.Open(receiver.getDSN(*receiver.connCredential)),
		&gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogLevel),
		},
	)
	if err != nil {
		logger.Error(err)

		return ErrConnectDB
	}

	db, err := gormDB.DB()
	if err != nil {
		logger.Error(err)

		return ErrConnectDB
	}

	db.SetMaxIdleConns(receiver.Config.MaxIdleConnections)
	db.SetMaxOpenConns(receiver.Config.MaxOpenConnections)
	db.SetConnMaxLifetime(time.Duration(receiver.Config.MaxLifetimeInMinutes) * time.Minute)

	var version string
	if res := gormDB.Raw("SELECT VERSION();").First(&version); res.Error != nil {
		logger.Error(res.Error)

		return ErrConnectDB
	}

	logger.Infof("Database version: %s", version)

	receiver.gorm = gormDB

	return nil
}
