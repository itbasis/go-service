package grpc

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
)

func InterceptorLogger(logger zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(
		func(ctx context.Context, level logging.Level, msg string, fields ...any) {
			logger = logger.With().Fields(fields).Logger()

			switch level {
			case logging.LevelDebug:
				logger.Debug().Msg(msg)
			case logging.LevelInfo:
				logger.Info().Msg(msg)
			case logging.LevelWarn:
				logger.Warn().Msg(msg)
			case logging.LevelError:
				logger.Error().Msg(msg)

			default:
				logger.Panic().Msgf("unknown level: %v", level)
			}
		},
	)
}
