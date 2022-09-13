package middleware

import (
	"bytes"
	"github.com/ferryvg/hiring-test-go-users-api/internal/transport/transportlib"

	"github.com/ferryvg/hiring-test-go-users-api/internal/dal"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var bearerAuthPrefix = []byte("Bearer ")

type JwtAuthChecker struct {
	logger  logrus.FieldLogger
	manager dal.UsersManager
}

func NewJwtAuthChecker(manager dal.UsersManager, logger logrus.FieldLogger) *JwtAuthChecker {
	return &JwtAuthChecker{
		manager: manager,
		logger:  logger,
	}
}

func (c *JwtAuthChecker) Middleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		auth := ctx.Request.Header.Peek("Authorization")
		if !bytes.HasPrefix(auth, bearerAuthPrefix) {
			c.unauthorizedError(ctx)
			return
		}

		payload := string(auth[len(bearerAuthPrefix):])

		user, err := c.manager.VerifyToken(ctx, string(payload))
		if err != nil {
			c.logger.WithField("token", string(payload)).Warn("Invalid or unknown access token")

			c.unauthorizedError(ctx)
			return
		}

		ctx.SetUserValue("user", user)

		next(ctx)
	}
}

func (c *JwtAuthChecker) unauthorizedError(ctx *fasthttp.RequestCtx) {
	status := fasthttp.StatusUnauthorized
	msg := fasthttp.StatusMessage(status)

	if err := transportlib.ResponseJsonStatus(ctx, msg, status); err != nil {
		ctx.Error(msg, status)
	}
}
