package service

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	// grpcZerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/itbasis/go-jwt-auth/grpc/client"
	"github.com/itbasis/go-jwt-auth/grpc/server"
	grpcLogUtils "github.com/itbasis/go-log-utils/grpc"
	itbasisServiceGrpc "github.com/itbasis/go-service/grpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Deprecated: use GetGrpcServer
func (receiver *Service) GetGrpc() *grpc.Server { return receiver.GetGrpcServer() }

func (receiver *Service) GetGrpcServer() *grpc.Server {
	if receiver.grpcServer != nil {
		return receiver.grpcServer
	}

	if receiver.config.GrpcServerDisabled {
		log.Error().Msg(gRPCServerIsDisabled)

		return nil
	}

	if receiver.grpcServer == nil {
		receiver.initGrpcServer()
	}

	return receiver.grpcServer
}

func (receiver *Service) GetGrpcServerMetrics() *grpcPrometheus.ServerMetrics {
	if receiver.grpcServerMetrics != nil {
		return receiver.grpcServerMetrics
	}

	if receiver.config.GrpcServerDisabled {
		log.Error().Msg(gRPCServerIsDisabled)

		return nil
	}

	if receiver.grpcServerMetrics == nil {
		receiver.InitGrpcServerMetrics(nil, nil)
	}

	return receiver.grpcServerMetrics
}

func (receiver *Service) initGrpcServer() {
	authFunc := server.NewAuthServerInterceptorWithCustomParser(receiver.jwtToken).GetAuthFunc()

	interceptorLogger := itbasisServiceGrpc.InterceptorLogger(log.Logger)

	var logOpts []logging.Option

	if log.Logger.Debug().Enabled() {
		logOpts = []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall, logging.PayloadReceived, logging.PayloadSent),
		}
	}

	unaryInterceptors := grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(interceptorLogger, logOpts...),
		grpcLogUtils.GrpcLogUnaryServerInterceptor(),
		auth.UnaryServerInterceptor(authFunc),
	)
	streamInterceptors := grpc.ChainStreamInterceptor(
		logging.StreamServerInterceptor(interceptorLogger, logOpts...),
		auth.StreamServerInterceptor(authFunc),
	)

	receiver.grpcServer = grpc.NewServer(unaryInterceptors, streamInterceptors)

	if receiver.config.GrpcReflectionEnabled {
		reflection.Register(receiver.grpcServer)
	}
}

func (receiver *Service) InitGrpcServerMetrics(
	serverMetricsOptions []grpcPrometheus.ServerMetricsOption,
	promHTTPHandlerOpts *promhttp.HandlerOpts,
) *Service {
	if len(serverMetricsOptions) == 0 {
		log.Debug().Msg("Using default server metrics...")

		serverMetricsOptions = []grpcPrometheus.ServerMetricsOption{
			grpcPrometheus.WithServerHandlingTimeHistogram(),
		}
	}

	if promHTTPHandlerOpts == nil {
		log.Debug().Msg("Using default server metric handlers...")

		promHTTPHandlerOpts = &promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}
	}

	log.Debug().Msgf("gRPC server metrics count: %d", len(serverMetricsOptions))

	serverMetrics := grpcPrometheus.NewServerMetrics(serverMetricsOptions...)
	registry := prometheus.NewPedanticRegistry()

	if err := registry.Register(serverMetrics); err != nil {
		log.Panic().Err(err).Send()
	}

	receiver.AddRestControllers(
		RestController{
			Method:  http.MethodGet,
			Path:    "/metrics/grpc",
			Handler: gin.WrapH(promhttp.HandlerFor(registry, *promHTTPHandlerOpts)),
		},
	)

	receiver.grpcServerMetrics = serverMetrics

	return receiver
}

func (receiver *Service) GetGrpcClientInterceptors(authClientInterceptor *client.AuthClientInterceptor) {
	interceptorLogger := itbasisServiceGrpc.InterceptorLogger(log.Logger)

	var logOpts []logging.Option
	if log.Logger.Debug().Enabled() {
		logOpts = []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall, logging.PayloadReceived, logging.PayloadSent),
		}
	}

	receiver.GrpcClientInterceptors = []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(interceptorLogger, logOpts...),
			authClientInterceptor.UnaryHeaderAuthorizeForwarder(),
		),

		grpc.WithChainStreamInterceptor(
			logging.StreamClientInterceptor(interceptorLogger, logOpts...),
			authClientInterceptor.UnaryStreamHeaderAuthorizeForwarder(),
		),
	}
}

func (receiver *Service) runGrpcServer(wg *sync.WaitGroup) {
	grpcServer := receiver.GetGrpcServer()

	if grpcServerMetrics := receiver.GetGrpcServerMetrics(); grpcServerMetrics != nil {
		grpcServerMetrics.InitializeMetrics(grpcServer)
	}

	if log.Debug().Enabled() {
		for service, info := range grpcServer.GetServiceInfo() {
			log.Debug().Msgf("service: %s , info: %+v\n", service, info)
		}
	}

	listen, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", receiver.config.GrpcServerHost, receiver.config.GrpcServerPort))
	if err != nil {
		log.Panic().Err(err).Send()
	}

	log.Info().Msgf("gRPC listen: %s", listen.Addr().String())

	if err = grpcServer.Serve(listen); err != nil {
		log.Panic().Err(err).Send()
	}

	wg.Done()
}
