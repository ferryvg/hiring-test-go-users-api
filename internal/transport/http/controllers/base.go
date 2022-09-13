package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/ferryvg/hiring-test-go-users-api/internal/transport/transportlib"

	"github.com/ferryvg/hiring-test-go-users-api/internal/domain"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type baseController struct {
	name   string
	logger logrus.FieldLogger
}

func newBaseController(name string, logger logrus.FieldLogger) baseController {
	return baseController{
		name:   name,
		logger: logger.WithField("controller", name),
	}
}

func (c *baseController) readPayload(ctx *fasthttp.RequestCtx, target interface{}) error {
	return json.Unmarshal(ctx.Request.Body(), target)
}

func (c *baseController) writePayload(ctx *fasthttp.RequestCtx, data interface{}) error {
	ctx.Response.ResetBody()

	payload, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to encode response body: %w", err)
	}

	ctx.Response.SetBody(payload)

	return nil
}

func (c *baseController) getUser(ctx *fasthttp.RequestCtx) *domain.User {
	user, ok := ctx.UserValue("user").(*domain.User)

	if !ok || user == nil {
		return domain.NewGuestUser()
	}

	return user
}

func (c *baseController) badRequestError(ctx *fasthttp.RequestCtx, msg string) {
	c.error(ctx, "", fasthttp.StatusBadRequest)
}

func (c *baseController) internalError(ctx *fasthttp.RequestCtx) {
	c.error(ctx, "", fasthttp.StatusInternalServerError)
}

func (c *baseController) error(ctx *fasthttp.RequestCtx, msg string, statusCode int) {
	if msg == "" {
		msg = fasthttp.StatusMessage(statusCode)
	}

	if err := transportlib.ResponseJsonStatus(ctx, msg, statusCode); err != nil {
		ctx.Error(msg, statusCode)
	}
}
