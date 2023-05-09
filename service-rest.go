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

type HttpController struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

// Deprecated: use HttpController
type RestController = HttpController

func (receiver *Service) GetGin() *gin.Engine {
	if receiver.gin == nil {
		receiver.initGinServer()
	}

	return receiver.gin
}

// Deprecated: use AddHttpControllers
func (receiver *Service) AddRestControllers(restControllers ...RestController) {
	receiver.AddHttpControllers(restControllers...)
}

func (receiver *Service) AddHttpControllers(httpControllers ...HttpController) {
	log.Debug().Msgf("REST controllers for adding: %v", httpControllers)

	receiver.ginControllers = append(receiver.ginControllers, httpControllers...)
	log.Debug().Msgf("REST controllers: %v", receiver.ginControllers)

	log.Trace().Msgf("gin: %v", receiver.gin)

	if receiver.gin != nil {
		for _, restController := range httpControllers {
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
