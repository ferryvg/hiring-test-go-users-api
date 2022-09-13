package controllers

import (
	"errors"

	"github.com/ferryvg/hiring-test-go-users-api/internal/dal"
	"github.com/ferryvg/hiring-test-go-users-api/internal/domain"
	"github.com/ferryvg/hiring-test-go-users-api/internal/transport/http/payloads"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type UsersController struct {
	baseController
	usersManager dal.UsersManager
	logger       logrus.FieldLogger
}

func NewUsersController(usersManager dal.UsersManager, logger logrus.FieldLogger) *UsersController {
	return &UsersController{
		usersManager:   usersManager,
		baseController: newBaseController("users", logger),
	}
}

func (c *UsersController) Me(ctx *fasthttp.RequestCtx) {
	user, ok := ctx.UserValue("user").(*domain.User)
	if !ok || user == nil {
		c.badRequestError(ctx, "")
		return
	}

	c.handleSingleUser(ctx, user)
}

func (c *UsersController) SingleUser(ctx *fasthttp.RequestCtx) {
	idUser, ok := ctx.UserValue("id_user").(string)
	if !ok || idUser == "" {
		c.badRequestError(ctx, "invalid user id")
		return
	}

	user, ok := ctx.UserValue("user").(*domain.User)
	if !ok || user == nil {
		c.badRequestError(ctx, "")
		return
	}

	if idUser == "me" {
		c.handleSingleUser(ctx, user)
		return
	}

	if !user.Roles[domain.AdminRole] {
		c.error(ctx, "", fasthttp.StatusForbidden)
		return
	}

	user, err := c.usersManager.Get(ctx, idUser)
	if err != nil || user == nil {
		if errors.Is(err, dal.ErrNotFound) {
			c.error(ctx, "", fasthttp.StatusNotFound)
			return
		}
	}

	c.handleSingleUser(ctx, user)
}

func (c *UsersController) ChangeRoles(ctx *fasthttp.RequestCtx) {
	idUser, ok := ctx.UserValue("id_user").(string)
	if !ok || idUser == "" {
		c.badRequestError(ctx, "invalid user id")
		return
	}

	var payload payloads.Roles
	if err := c.readPayload(ctx, &payload); err != nil {
		c.badRequestError(ctx, "invalid roles payload")
		return
	}

	roles := make(map[domain.Role]bool, len(payload.Roles))
rolesLoop:
	for roleStr, enabled := range payload.Roles {
		role := domain.RoleFromString(roleStr)

		switch role {
		case domain.GuestRole, domain.UnknownRole:
			continue rolesLoop
		}

		roles[role] = enabled
	}

	if err := c.usersManager.ChangeRoles(ctx, idUser, roles); err != nil {
		c.logger.WithError(err).WithField("user", idUser).WithField("roles", roles).Error("Failed to update user roles")
		c.internalError(ctx)
		return
	}

	c.error(ctx, "", fasthttp.StatusNoContent)
}

func (c *UsersController) UsersList(ctx *fasthttp.RequestCtx) {
	usersList, err := c.usersManager.GetList(ctx)
	if err != nil {
		c.logger.WithError(err).Error("Failed to load users list")
		c.internalError(ctx)
		return
	}

	payload := &payloads.UsersList{
		Users: make([]payloads.User, 0, len(usersList)),
	}

	for _, user := range usersList {
		payload.Users = append(payload.Users, payloads.NewUser(&user))
	}

	if err := c.writePayload(ctx, payload); err != nil {
		c.logger.WithError(err).Error("Failed to write users list to response")
		c.internalError(ctx)
	}
}

func (c *UsersController) handleSingleUser(ctx *fasthttp.RequestCtx, user *domain.User) {
	payload := payloads.NewUser(user)

	if err := c.writePayload(ctx, payload); err != nil {
		c.logger.WithError(err).Error("Failed to write response payload")
		c.internalError(ctx)
	}
}
