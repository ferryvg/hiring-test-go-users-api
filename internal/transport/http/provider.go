package http

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core/http"
	coreMiddleware "github.com/ferryvg/hiring-test-go-users-api/internal/core/http/middleware"
	"github.com/ferryvg/hiring-test-go-users-api/internal/dal"
	"github.com/ferryvg/hiring-test-go-users-api/internal/domain"
	"github.com/ferryvg/hiring-test-go-users-api/internal/transport/http/controllers"
	"github.com/ferryvg/hiring-test-go-users-api/internal/transport/http/middleware"
	"github.com/sirupsen/logrus"
)

type Provider struct{}

func (p *Provider) Register(c core.Container) {
	p.registerMiddlewares(c)
	p.registerAuthController(c)
	p.registerUsersController(c)

	p.registerUsersEndpoints(c)
}

func (*Provider) registerMiddlewares(c core.Container) {
	c.Set("svc.transport.http.middleware.jwt", func(c core.Container) interface{} {
		manager := c.MustGet("svc.dal.users").(dal.UsersManager)
		logger := c.MustGet("logger").(logrus.FieldLogger)

		return middleware.NewJwtAuthChecker(manager, logger)
	})

	c.Set("svc.transport.http.middleware.role", func(c core.Container) interface{} {
		logger := c.MustGet("logger").(logrus.FieldLogger)

		return middleware.NewRoleChecker(logger)
	})
}

func (*Provider) registerAuthController(c core.Container) {
	c.Set("svc.transport.http.controllers.auth", func(c core.Container) interface{} {
		manager := c.MustGet("svc.dal.users").(dal.UsersManager)
		logger := c.MustGet("logger").(logrus.FieldLogger)

		return controllers.NewAuthController(manager, logger)
	})

	c.MustExtend("http.router", func(old interface{}, c core.Container) interface{} {
		router := old.(*fasthttprouter.Router)
		logger := c.MustGet("logger").(logrus.FieldLogger)

		controller := c.MustGet("svc.transport.http.controllers.auth").(*controllers.AuthController)
		roleChecker := c.MustGet("svc.transport.http.middleware.role").(*middleware.RoleChecker)

		middlewares := []http.Middleware{
			roleChecker.Middleware([]domain.Role{domain.GuestRole}),
			coreMiddleware.Metrics,
			coreMiddleware.Logger(logger),
		}

		signInEndpoint := http.BuildHandler(
			controller.SignIn,
			middlewares...,
		)

		signUpEndpoint := http.BuildHandler(
			controller.SignUp,
			middlewares...,
		)

		router.POST("/login", signInEndpoint)
		router.POST("/register", signUpEndpoint)

		return router
	})
}

func (*Provider) registerUsersController(c core.Container) {
	c.Set("svc.transport.http.controllers.users", func(c core.Container) interface{} {
		manager := c.MustGet("svc.dal.users").(dal.UsersManager)
		logger := c.MustGet("logger").(logrus.FieldLogger)

		return controllers.NewUsersController(manager, logger)
	})
}

func (*Provider) registerUsersEndpoints(c core.Container) {
	c.MustExtend("http.router", func(old interface{}, c core.Container) interface{} {
		router := old.(*fasthttprouter.Router)
		logger := c.MustGet("logger").(logrus.FieldLogger)

		controller := c.MustGet("svc.transport.http.controllers.users").(*controllers.UsersController)
		jwtChecker := c.MustGet("svc.transport.http.middleware.jwt").(*middleware.JwtAuthChecker)
		roleChecker := c.MustGet("svc.transport.http.middleware.role").(*middleware.RoleChecker)

		adminMiddlewares := []http.Middleware{
			roleChecker.Middleware([]domain.Role{domain.AdminRole}),
			jwtChecker.Middleware,
			coreMiddleware.Metrics,
			coreMiddleware.Logger(logger),
		}

		userEndpoint := http.BuildHandler(
			controller.SingleUser,
			roleChecker.Middleware([]domain.Role{domain.BasicRole, domain.AdminRole}),
			jwtChecker.Middleware,
			coreMiddleware.Metrics,
			coreMiddleware.Logger(logger),
		)
		usersListEndpoint := http.BuildHandler(controller.UsersList, adminMiddlewares...)
		changeRolesEndpoint := http.BuildHandler(controller.ChangeRoles, adminMiddlewares...)

		router.GET("/users/:id_user", userEndpoint)
		router.PUT("/users/:id_user", changeRolesEndpoint)
		router.GET("/users/", usersListEndpoint)

		return router
	})
}
