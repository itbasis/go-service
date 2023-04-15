package service

import (
	"sync"

	"github.com/benbjohnson/clock"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	coreUtils "github.com/itbasis/go-core-utils"
	jwtToken "github.com/itbasis/go-jwt-auth/jwt-token"
	logUtils "github.com/itbasis/go-log-utils"
	"github.com/itbasis/go-service/db"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Service struct {
	config *Config

	jwtToken jwtToken.JwtToken

	grpcServer             *grpc.Server
	grpcServerMetrics      *prometheus.ServerMetrics
	GrpcClientInterceptors []grpc.DialOption

	gin            *gin.Engine
	ginControllers []RestController

	clock     clock.Clock
	gorm      *db.DB
	scheduler *gocron.Scheduler
}

func NewServiceWithEnvironment() *Service {
	logUtils.ConfigureDefaultContextLogger(false)

	config := &Config{}
	if err := coreUtils.ReadEnvConfig(config); err != nil {
		log.Panic().Err(err).Msg("")
	}

	return NewServiceWithConfig(config)
}

func NewServiceWithConfig(config *Config) *Service {
	log.Logger.Debug().Msgf("config: %++v", config)

	_, err := logUtils.ConfigureRootLogger(config.ServiceName)
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	service := &Service{config: config}
	service.clock = clock.New()
	service.initJwtToken()

	return service
}

func (receiver *Service) Run() {
	log.Debug().Msg("running service...")

	if receiver.config.SchedulerEnabled {
		receiver.scheduler.StartAsync()
	}

	wg := &sync.WaitGroup{}

	// gRPC
	if receiver.config.GrpcServerDisabled {
		log.Info().Msg(gRPCServerIsDisabled)
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
