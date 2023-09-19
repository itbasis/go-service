package service

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/itbasis/go-clock"
	coreUtils "github.com/itbasis/go-core-utils/v2"
	jwtToken "github.com/itbasis/go-jwt-auth/v2/jwt-token"
	logUtils "github.com/itbasis/go-log-utils/v2"
	"github.com/itbasis/go-service/v2/db"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Service struct {
	config *Config

	jwtToken jwtToken.JwtToken

	grpcServer             *grpc.Server
	grpcServerMetrics      *prometheus.ServerMetrics
	GrpcClientInterceptors []grpc.DialOption

	gin            *gin.Engine
	ginControllers []HTTPController

	clock     clock.Clock
	gorm      *db.DB
	scheduler *gocron.Scheduler
}

func NewServiceWithEnvironment(ctx context.Context, zapConfig zap.Config) *Service {
	logUtils.ConfigureDefaultContextLogger(false, zapConfig)

	config := &Config{}
	if err := coreUtils.ReadEnvConfig(ctx, config, nil); err != nil {
		zapctx.Default.Sugar().Panic(err)
	}

	return NewServiceWithConfig(ctx, zapConfig, config)
}

func NewServiceWithConfig(ctx context.Context, zapConfig zap.Config, config *Config) *Service {
	logger := zapctx.Logger(ctx).Sugar()
	logger.Debugf("config: %++v", config)

	_, err := logUtils.ConfigureRootLogger(ctx, config.ServiceName, zapConfig)
	if err != nil {
		logger.Panic(err)
	}

	service := &Service{config: config}
	service.clock = clock.New()
	service.initJwtToken(ctx)

	return service
}

func (receiver *Service) Run() {
	logger := zapctx.Default.Sugar()
	logger.Debug("running service...")

	if receiver.config.SchedulerEnabled {
		receiver.scheduler.StartAsync()
	}

	wg := &sync.WaitGroup{}

	// gRPC
	if receiver.config.GrpcServerDisabled {
		logger.Info(gRPCServerIsDisabled)
	} else {
		wg.Add(1)
		go receiver.runGrpcServer(wg)
	}

	// HTTP
	wg.Add(1)

	go receiver.runGinServer(wg)

	wg.Wait()
}

func (receiver *Service) GetConfig() Config {
	return *receiver.config
}

func (receiver *Service) GetClock() clock.Clock {
	return receiver.clock
}
