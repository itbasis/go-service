package service

import (
	"fmt"
	"sync"

	ginZerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/itbasis/go-service/rest"
	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
)

type RestController struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

func (receiver *Service) GetGin() *gin.Engine {
	if receiver.gin == nil {
		receiver.initGinServer()
	}

	return receiver.gin
}

func (receiver *Service) AddRestControllers(restControllers ...RestController) {
	log.Debug().Msgf("REST controllers for adding: %v", restControllers)

	receiver.ginControllers = append(receiver.ginControllers, restControllers...)
	log.Debug().Msgf("REST controllers: %v", receiver.ginControllers)

	log.Trace().Msgf("gin: %v", receiver.gin)

	if receiver.gin != nil {
		for _, restController := range restControllers {
			log.Debug().Msgf("adding REST controller: %v", restController)

			receiver.gin.Handle(restController.Method, restController.Path, restController.Handler)
		}
	}
}

func (receiver *Service) initGinServer() {
	receiver.gin = gin.New()
	receiver.gin.Use(
		gin.Recovery(),
		ginZerolog.Logger("rest"),
		receiver.ginLoggerCtx,
	)

	if log.Trace().Enabled() {
		log.Info().Msg("enable REST request tracing")

		receiver.gin.Use(rest.LoggingRequest)
	}

	log.Debug().Msgf("REST controllers: %v", receiver.ginControllers)

	for _, restController := range receiver.ginControllers {
		log.Debug().Msgf("adding REST controller: %v", restController)

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
		log.Panic().Err(err).Send()
	}

	wg.Done()
}
