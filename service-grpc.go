package service

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	// grpcZerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/itbasis/go-jwt-auth/grpc/client"
	"github.com/itbasis/go-jwt-auth/grpc/server"
	grpcLogUtils "github.com/itbasis/go-log-utils/grpc"
	itbasisServiceGrpc "github.com/itbasis/go-service/grpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (receiver *Service) GetGrpc() *grpc.Server {
	if receiver.config.GrpcServerDisabled {
		log.Error().Msg(gRPCServerIsDisabled)

		return nil
	}

	if receiver.grpc == nil {
		receiver.initGrpcServer()
	}

	return receiver.grpc
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
		grpcPrometheus.UnaryServerInterceptor,
		grpcLogUtils.GrpcLogUnaryServerInterceptor(),
		auth.UnaryServerInterceptor(authFunc),
	)
	streamInterceptors := grpc.ChainStreamInterceptor(
		logging.StreamServerInterceptor(interceptorLogger, logOpts...),
		grpcPrometheus.StreamServerInterceptor,
		auth.StreamServerInterceptor(authFunc),
	)

	receiver.grpc = grpc.NewServer(unaryInterceptors, streamInterceptors)

	if receiver.config.GrpcReflectionEnabled {
		reflection.Register(receiver.grpc)
	}
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
	grpcServer := receiver.GetGrpc()

	grpcPrometheus.Register(grpcServer)
	http.Handle("/metrics", promhttp.Handler())

	if log.Debug().Enabled() {
		for service, info := range grpcServer.GetServiceInfo() {
			log.Debug().Msgf("service: %s , info: %+v\n", service, info)
		}
	}

	listen, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", receiver.config.GrpcServerHost, receiver.config.GrpcServerPort))
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	log.Info().Msgf("gRPC listen: %s", listen.Addr().String())

	if err = grpcServer.Serve(listen); err != nil {
		log.Panic().Err(err).Msg("")
	}

	wg.Done()
}
