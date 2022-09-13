package middleware

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/domain"
	"github.com/ferryvg/hiring-test-go-users-api/internal/transport/transportlib"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type RoleChecker struct {
	logger logrus.FieldLogger
}

func NewRoleChecker(logger logrus.FieldLogger) *RoleChecker {
	return &RoleChecker{
		logger: logger,
	}
}

func (c *RoleChecker) Middleware(roles []domain.Role) func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			user, ok := ctx.UserValue("user").(*domain.User)
			if !ok || user == nil {
				c.logger.WithField("request_uri", ctx.Request.URI().String()).Warn("Unauthorized request")

				user = domain.NewGuestUser()
			}

			for _, role := range roles {
				if user.Roles[role] {
					next(ctx)
					return
				}
			}

			status := fasthttp.StatusForbidden
			msg := fasthttp.StatusMessage(fasthttp.StatusForbidden)

			if err := transportlib.ResponseJsonStatus(ctx, msg, status); err != nil {
				ctx.Error(msg, status)
			}
		}
	}
}
