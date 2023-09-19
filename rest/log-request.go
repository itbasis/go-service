package rest

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap/zapcore"
)

func LoggingRequest(ctx *gin.Context) {
	log := zapctx.Logger(ctx)
	logger := log.Sugar()

	var buf bytes.Buffer

	tee := io.TeeReader(ctx.Request.Body, &buf)
	body, err := io.ReadAll(tee)

	if err != nil {
		logger.Error(err)

		ctx.Next()

		return
	}

	ctx.Request.Body = io.NopCloser(&buf)

	if log.Core().Enabled(zapcore.DebugLevel) {
		logger.Debugf("uri: %s", ctx.Request.RequestURI)
		logger.Debugf("headers: %++v", ctx.Request.Header)
		logger.Debugf("body: %s", string(body))
	}

	ctx.Next()
}
