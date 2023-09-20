package service

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/itbasis/go-jwt-auth/v2/grpc/client"
	"github.com/itbasis/go-jwt-auth/v2/grpc/server"
	grpcLogUtils "github.com/itbasis/go-log-utils/v2/grpc"
	itbasisServiceGrpc "github.com/itbasis/go-service/v2/grpc"
	"github.com/juju/zaputil/zapctx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		zapctx.Default.Sugar().Error(gRPCServerIsDisabled)

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
		zapctx.Default.Sugar().Error(gRPCServerIsDisabled)

		return nil
	}

	if receiver.grpcServerMetrics == nil {
		receiver.InitGrpcServerMetrics(nil, nil)
	}

	return receiver.grpcServerMetrics
}

func (receiver *Service) initGrpcServer() {
	authFunc := server.NewAuthServerInterceptorWithCustomParser(receiver.jwtToken).GetAuthFunc()

	logger := zapctx.Default
	interceptorLogger := itbasisServiceGrpc.InterceptorLogger(logger)

	var logOpts []logging.Option

	if logger.Core().Enabled(zapcore.DebugLevel) {
		logOpts = []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall, logging.PayloadReceived, logging.PayloadSent),
		}
	}

	unaryInterceptors := grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(interceptorLogger, logOpts...),
		grpcLogUtils.LogUnaryServerInterceptor(),
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
	log := zapctx.Default.Sugar()

	if len(serverMetricsOptions) == 0 {
		log.Debug("Using default server metrics...")

		serverMetricsOptions = []grpcPrometheus.ServerMetricsOption{
			grpcPrometheus.WithServerHandlingTimeHistogram(),
		}
	}

	if promHTTPHandlerOpts == nil {
		log.Debug("Using default server metric handlers...")

		promHTTPHandlerOpts = &promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}
	}

	log.Debugf("gRPC server metrics count: %d", len(serverMetricsOptions))

	serverMetrics := grpcPrometheus.NewServerMetrics(serverMetricsOptions...)
	registry := prometheus.NewPedanticRegistry()

	if err := registry.Register(serverMetrics); err != nil {
		log.Panic(err)
	}

	receiver.AddHTTPControllers(
		HTTPController{
			Method:  http.MethodGet,
			Path:    "/metrics/grpc",
			Handler: gin.WrapH(promhttp.HandlerFor(registry, *promHTTPHandlerOpts)),
		},
	)

	receiver.grpcServerMetrics = serverMetrics

	return receiver
}

func (receiver *Service) GetGrpcClientInterceptors(authClientInterceptor *client.AuthClientInterceptor) {
	logger := zapctx.Default

	interceptorLogger := itbasisServiceGrpc.InterceptorLogger(logger)

	var logOpts []logging.Option
	if logger.Core().Enabled(zapcore.DebugLevel) {
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

	logger := zapctx.Default
	log := logger.Sugar()

	if grpcServerMetrics := receiver.GetGrpcServerMetrics(); grpcServerMetrics != nil {
		grpcServerMetrics.InitializeMetrics(grpcServer)
	}

	if logger.Core().Enabled(zapcore.DebugLevel) {
		for service, info := range grpcServer.GetServiceInfo() {
			log.Debugf("service: %s , info: %+v\n", service, info)
		}
	}

	listen, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", receiver.config.GrpcServerHost, receiver.config.GrpcServerPort))
	if err != nil {
		log.Error(err)
		log.Panic(fmt.Errorf(msgErrFailedStartGRPCServer, err))
	}

	log.Infof("gRPC listen: %s", listen.Addr().String())

	if err = grpcServer.Serve(listen); err != nil {
		log.Error(err)
		log.Panic(fmt.Errorf(msgErrFailedStartGRPCServer, err))
	}

	wg.Done()
}
