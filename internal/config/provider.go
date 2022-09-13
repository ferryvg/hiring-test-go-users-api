package config

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/sirupsen/logrus"
)

type Provider struct{}

func (p *Provider) Register(c core.Container) {
	p.registerBuilder(c)
	p.registerConfig(c)
}

func (*Provider) registerBuilder(c core.Container) {
	c.Set("svc.config.builder", func(c core.Container) interface{} {
		logger := c.MustGet("logger").(logrus.FieldLogger)

		return NewBuilder(logger)
	})
}

func (*Provider) registerConfig(c core.Container) {
	c.Set("svc.config", func(c core.Container) interface{} {
		builder := c.MustGet("svc.config.builder").(Builder)
		confFile := c.MustGet("svc.config_file").(*string)
		logger := c.MustGet("logger").(logrus.FieldLogger)

		config, err := builder.Build(*confFile)
		if err != nil {
			logger.WithError(err).Fatal("Failed to parse service configuration")
		}

		config.Mysql.Nodes = NewClusterNodeStore()

		return config
	})
}
