package service

import (
	jwtToken "github.com/itbasis/go-jwt-auth/jwt-token"
	jwtAuthTokenImpl "github.com/itbasis/go-jwt-auth/jwt-token/impl"
	"github.com/rs/zerolog/log"
)

func (receiver *Service) initJwtToken() {
	jt, err := jwtAuthTokenImpl.NewJwtToken(receiver.clock)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	receiver.jwtToken = jt
}

func (receiver *Service) GetJwtToken() jwtToken.JwtToken { return receiver.jwtToken }
