package grpc

import (
	"context"
	"crypto/tls"

	"github.com/juju/zaputil/zapctx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func GetServiceConnection(ctx context.Context, serviceHost string, useSSL bool, opts ...grpc.DialOption) *grpc.ClientConn {
	logger := zapctx.Logger(ctx).Sugar()

	logger.Debugf("getting service connection for host: %s", serviceHost)

	if useSSL {
		// FIXME Add SSL connection method with certificate verification
		logger.Warnf("The connection will be with an insecure SSL connection (InsecureSkipVerify=true) for the host: %s", serviceHost)

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
		logger.Panic(err)
	}

	logger.Infof("connection state for host '%s': %s", serviceHost, conn.GetState())

	return conn
}
