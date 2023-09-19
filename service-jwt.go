package service

import (
	"context"

	jwtToken "github.com/itbasis/go-jwt-auth/v2/jwt-token"
	jwtAuthTokenImpl "github.com/itbasis/go-jwt-auth/v2/jwt-token/impl"
	"github.com/juju/zaputil/zapctx"
)

func (receiver *Service) initJwtToken(ctx context.Context) {
	jt, err := jwtAuthTokenImpl.NewJwtToken(receiver.clock)
	if err != nil {
		zapctx.Logger(ctx).Panic(err.Error())
	}

	receiver.jwtToken = jt
}

func (receiver *Service) GetJwtToken() jwtToken.JwtToken { return receiver.jwtToken }
