package middleware

import (
	"io"
	"time"

	"github.com/go-kit/log"
	"github.com/zackattackz/statikit-render-svc/internal/models"
	"github.com/zackattackz/statikit-render-svc/internal/service"
)

func LoggingMW(logger log.Logger) service.RenderServiceMiddleware {
	return func(next service.RenderService) service.RenderService {
		return logMW{logger, next}
	}
}

type logMW struct {
	logger log.Logger
	service.RenderService
}

func (mw logMW) Render(contents string, schema models.Schema, w io.Writer) error {
	var err error
	k := service.CacheKey{Schema: schema, Contents: contents}
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "render",
			"cache-key", k.Hash(),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.RenderService.Render(contents, schema, w)
	return err
}
