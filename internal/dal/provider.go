package dal

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/config"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/ferryvg/hiring-test-go-users-api/internal/db"
	"github.com/sirupsen/logrus"
)

type Provider struct{}

func (p *Provider) Register(c core.Container) {
	p.registerUsersManager(c)
}

func (*Provider) registerUsersManager(c core.Container) {
	c.Set("svc.dal.users", func(c core.Container) interface{} {
		conf := c.MustGet("svc.config").(*config.AppConfig)
		dbManager := c.MustGet("svc.db.manager").(db.Manager)
		logger := c.MustGet("logger").(logrus.FieldLogger)

		return NewUsersManager(&conf.Jwt, dbManager, logger)
	})
}
