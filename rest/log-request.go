package rest

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LoggingRequest(ctx *gin.Context) {
	logger := zerolog.Ctx(ctx)

	var buf bytes.Buffer

	tee := io.TeeReader(ctx.Request.Body, &buf)
	body, err := io.ReadAll(tee)

	if err != nil {
		logger.Error().Err(err).Send()

		ctx.Next()

		return
	}

	ctx.Request.Body = io.NopCloser(&buf)

	log.Trace().Msgf("uri: %s", ctx.Request.RequestURI)
	log.Trace().Msgf("headers: %++v", ctx.Request.Header)
	log.Trace().Msgf("body: %s", string(body))

	ctx.Next()
}
