package service

import "google.golang.org/protobuf/types/known/wrapperspb"

var (
	ReturnFalse = wrapperspb.Bool(false)
	ReturnTrue  = wrapperspb.Bool(true)
)

const (
	gormIsDisabled       = "GORM is disabled"
	httpServerIsDisabled = "HTTP server is disabled"
	gRPCServerIsDisabled = "gRPC server is disabled"
)
