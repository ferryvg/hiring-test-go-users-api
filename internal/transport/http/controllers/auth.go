package controllers

import (
	"errors"

	"github.com/ferryvg/hiring-test-go-users-api/internal/domain"

	"github.com/ferryvg/hiring-test-go-users-api/internal/dal"
	"github.com/ferryvg/hiring-test-go-users-api/internal/transport/http/payloads"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type AuthController struct {
	baseController
	usersManager dal.UsersManager
}

func NewAuthController(usersManager dal.UsersManager, logger logrus.FieldLogger) *AuthController {
	return &AuthController{
		baseController: newBaseController("auth", logger),
		usersManager:   usersManager,
	}
}

func (c *AuthController) SignIn(ctx *fasthttp.RequestCtx) {
	var authPayload payloads.Authenticate
	if err := c.readPayload(ctx, &authPayload); err != nil {
		c.badRequestError(ctx, "invalid credentials")
		return
	}

	token, err := c.usersManager.Authenticate(ctx, authPayload.ID, authPayload.Secret)
	if err != nil {
		if errors.Is(err, dal.ErrInvalidCredentials) {
			c.badRequestError(ctx, "invalid credentials")
			return
		}

		c.logger.WithError(err).Error("Failed to check user credentials")
		c.internalError(ctx)
		return
	}

	payload := payloads.Jwt{
		AccessToken: token.AccessToken,
		ExpiredAt:   token.ExpiredAt,
	}

	if err := c.writePayload(ctx, payload); err != nil {
		c.logger.WithError(err).Error("Failed to write access token to response")
		c.internalError(ctx)
	}
}

func (c *AuthController) SignUp(ctx *fasthttp.RequestCtx) {
	var userPayload payloads.SignUp
	if err := c.readPayload(ctx, &userPayload); err != nil {
		c.badRequestError(ctx, "invalid request data")
		return
	}

	user := &domain.User{
		ID:     userPayload.ID,
		Secret: userPayload.Secret,
		Roles: map[domain.Role]bool{
			domain.BasicRole: true,
		},
	}

	// TODO: prevent create users with registered IDs like a 'me', 'guest', etc.

	if err := c.usersManager.Create(ctx, user); err != nil {
		if errors.Is(err, dal.ErrUserExists) {
			c.badRequestError(ctx, "user with provided id already exists")
			return
		}

		c.logger.WithError(err).Error("Failed to register user")
		c.internalError(ctx)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
}
