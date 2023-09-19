package grpc

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

func InterceptorLogger(logger *zap.Logger) logging.Logger { //nolint:cyclop
	return logging.LoggerFunc(
		func(ctx context.Context, level logging.Level, msg string, fields ...any) {
			//nolint:gomnd
			zapFields := make([]zap.Field, 0, len(fields)/2)

			for i := 0; i < len(fields); i += 2 {
				key := fields[i].(string)
				value := fields[i+1]

				switch vType := value.(type) {
				case string:
					zapFields = append(zapFields, zap.String(key, vType))
				case int:
					zapFields = append(zapFields, zap.Int(key, vType))
				case bool:
					zapFields = append(zapFields, zap.Bool(key, vType))
				default:
					zapFields = append(zapFields, zap.Any(key, vType))
				}
			}

			log := logger.With(zapFields...)

			switch level {
			case logging.LevelDebug:
				log.Debug(msg)
			case logging.LevelInfo:
				log.Info(msg)
			case logging.LevelWarn:
				log.Warn(msg)
			case logging.LevelError:
				log.Error(msg)

			default:
				log.Panic(fmt.Sprintf("unknown level: %v", level))
			}
		},
	)
}
