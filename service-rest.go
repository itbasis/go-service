package service

import (
	"fmt"
	"sync"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/itbasis/go-service/v2/rest"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap/zapcore"
)

type HTTPController struct {
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

func (receiver *Service) AddHTTPControllers(httpControllers ...HTTPController) {
	log := zapctx.Default.Sugar()
	log.Debugf("REST controllers for adding: %v", httpControllers)

	receiver.ginControllers = append(receiver.ginControllers, httpControllers...)
	log.Debugf("REST controllers: %v", receiver.ginControllers)

	log.Debugf("gin: %v", receiver.gin)

	if receiver.gin != nil {
		for _, restController := range httpControllers {
			log.Debugf("adding REST controller: %v", restController)

			receiver.gin.Handle(restController.Method, restController.Path, restController.Handler)
		}
	}
}

func (receiver *Service) initGinServer() {
	log := zapctx.Default

	receiver.gin = gin.New()
	receiver.gin.Use(
		gin.Recovery(),
		ginzap.RecoveryWithZap(zapctx.Default, true),
	)

	if log.Core().Enabled(zapcore.DebugLevel) {
		log.Info("enable REST request tracing")

		receiver.gin.Use(rest.LoggingRequest)
	}

	log.Debug(fmt.Sprintf("REST controllers: %v", receiver.ginControllers))

	for _, restController := range receiver.ginControllers {
		log.Debug(fmt.Sprintf("adding REST controller: %v", restController))

		receiver.gin.Handle(restController.Method, restController.Path, restController.Handler)
	}
}

func (receiver *Service) runGinServer(wg *sync.WaitGroup) {
	log := zapctx.Default

	engine := receiver.GetGin()

	if engine == nil {
		wg.Done()

		return
	}

	listen := fmt.Sprintf("%s:%d", receiver.config.RestServerHost, receiver.config.RestServerPort)
	log.Debug(fmt.Sprintf("rest listen address: %s", listen))

	if err := engine.Run(listen); err != nil {
		log.Panic(err.Error())
	}

	wg.Done()
}
