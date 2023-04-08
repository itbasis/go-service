package grpc

import (
	"crypto/tls"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func GetServiceConnection(logger zerolog.Logger, serviceHost string, useSSL bool, opts ...grpc.DialOption) *grpc.ClientConn {
	logger.Debug().Msgf("getting service connection for host: %s", serviceHost)

	if useSSL {
		// FIXME Add SSL connection method with certificate verification
		logger.Warn().Msgf("The connection will be with an insecure SSL connection (InsecureSkipVerify=true) for the host: %s", serviceHost)

		// #nosec
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}))) //nolint:gosec
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(
		serviceHost,
		opts...,
	)
	if nil != err {
		logger.Panic().Err(err).Msg("")
	}

	logger.Info().Msgf("connection state for host '%s': %s", serviceHost, conn.GetState())

	return conn
}
