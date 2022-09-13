package middleware

import (
	"time"

	"github.com/ferryvg/hiring-test-go-users-api/internal/core/http"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func Logger(logger logrus.FieldLogger) http.Middleware {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			startedAt := time.Now()

			next(ctx)

			fields := logrus.Fields{
				"http.code":     ctx.Response.StatusCode(),
				"http.time_ms":  int64(time.Since(startedAt) / time.Millisecond),
				"http.method":   string(ctx.Request.Header.Method()),
				"http.endpoint": string(ctx.Request.URI().Path()),
			}

			getResponse := func(ctx *fasthttp.RequestCtx) string {
				body := ""
				if len(ctx.Response.Body()) > 1024 {
					body = string(ctx.Response.Body()[:1024])
				} else {
					body = string(ctx.Response.Body())
				}

				return body
			}

			debugFields := func(ctx *fasthttp.RequestCtx) logrus.Fields {
				return logrus.Fields{
					"ip":       string(ctx.Request.Header.Peek("X-Real-IP")),
					"refferer": string(ctx.Referer()),
					"ua":       string(ctx.UserAgent()),
					"uri":      ctx.Request.URI().String(),
				}
			}

			if ctx.Response.StatusCode() < 400 {
				logger.WithFields(fields).Info("finished call")
			} else if ctx.Response.StatusCode() < 500 {
				logger.WithFields(fields).WithFields(debugFields(ctx)).Warn("finished call: " + getResponse(ctx))
			} else {
				logger.WithFields(fields).WithFields(debugFields(ctx)).Error("finished call: " + getResponse(ctx))
			}
		}
	}
}
