package service

import (
	"fmt"
	"sync"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	// "github.com/pereslava/grpc_zerolog/ctxzerolog"
	"github.com/rs/zerolog/log"
)

type RestController struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

func (receiver *Service) GetGin() *gin.Engine {
	if receiver.config.RestServerDisabled {
		log.Error().Msg(gRPCServerIsDisabled)

		return nil
	}

	if receiver.gin == nil {
		receiver.initGinServer()
	}

	return receiver.gin
}

func (receiver *Service) AddRestControllers(restControllers ...RestController) {
	log.Trace().Msgf("REST controllers for adding: %v", restControllers)

	receiver.ginControllers = append(receiver.ginControllers, restControllers...)
	log.Trace().Msgf("REST controllers: %v", receiver.ginControllers)
}

func (receiver *Service) initGinServer() {
	receiver.gin = gin.New()
	receiver.gin.Use(
		gin.Recovery(),
		receiver.ginLoggerCtx,
		ginzerolog.Logger("rest"),
	)

	log.Debug().Msgf("REST controllers: %v", receiver.ginControllers)

	for _, restController := range receiver.ginControllers {
		receiver.gin.Handle(restController.Method, restController.Path, restController.Handler)
	}
}

func (receiver *Service) ginLoggerCtx(ctx *gin.Context) {
	log.Trace().Msg("Setting Logger in context")

	zerolog.Ctx(ctx)
}

func (receiver *Service) runGinServer(wg *sync.WaitGroup) {
	engine := receiver.GetGin()

	if engine == nil {
		wg.Done()

		return
	}

	listen := fmt.Sprintf("%s:%d", receiver.config.RestServerHost, receiver.config.RestServerPort)
	log.Trace().Msgf("rest listen address: %s", listen)

	if err := engine.Run(listen); err != nil {
		log.Panic().Err(err).Msg("")
	}

	wg.Done()
}
