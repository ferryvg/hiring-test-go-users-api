package db

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/config"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core/sd"
	"github.com/sirupsen/logrus"
)

type Provider struct{}

// Register components that provides access to MySQL database.
func (p *Provider) Register(c core.Container) {
	p.registerMysqlFactory(c)
	p.registerMysqlConnList(c)
	p.registerMysqlResolver(c)
	p.registerMysqlManager(c)
}

// Boot initialize database connection manager.
func (*Provider) Boot(c core.Container) error {
	logger := c.MustGet("logger").(logrus.FieldLogger)
	manager := c.MustGet("svc.db.manager").(Manager)

	logger.Info("Initialize MySQL manager")
	if err := manager.Init(); err != nil {
		logger.WithError(err).Error("Failed to initialize MySQL manager")

		return err
	}
	logger.Info("Mysql manager initialized successfully")

	return nil
}

// Shutdown database connection manager.
func (*Provider) Shutdown(c core.Container) {
	manager := c.MustGet("svc.db.manager").(Manager)
	manager.Shutdown()
}

// Register connection factory.
func (*Provider) registerMysqlFactory(c core.Container) {
	c.Set("svc.db.conn_factory", func(c core.Container) interface{} {
		conf := c.MustGet("svc.config").(*config.AppConfig)

		return NewMysqlConnFactory(conf.Mysql.Database, conf.Mysql.Username, conf.Mysql.Password)
	})
}

// Register opened connections list.
func (*Provider) registerMysqlConnList(c core.Container) {
	c.Set("svc.db.conn_list", func(c core.Container) interface{} {
		factory := c.MustGet("svc.db.conn_factory").(ConnFactory)
		logger := c.MustGet("logger").(logrus.FieldLogger)
		conf := c.MustGet("svc.config").(*config.AppConfig)

		mysqlConf := &conf.Mysql

		return NewConnList(factory, mysqlConf.Nodes, logger)
	})
}

// Register mysql nodes resolver.
func (*Provider) registerMysqlResolver(c core.Container) {
	c.Set("svc.db.resolver", func(c core.Container) interface{} {
		registry := c.MustGet("sd.registry").(sd.Registry)
		connList := c.MustGet("svc.db.conn_list").(ConnList)
		logger := c.MustGet("logger").(logrus.FieldLogger)
		conf := &ResolverConf{
			Service: "mysql",
			Tags:    []string{c.MustGet("saas.cluster").(string)},
		}

		return NewResolver(registry, connList, logger, conf)
	})
}

// Register mysql connections manager.
func (*Provider) registerMysqlManager(c core.Container) {
	c.Set("svc.db.manager", func(c core.Container) interface{} {
		connList := c.MustGet("svc.db.conn_list").(ConnList)
		resolver := c.MustGet("svc.db.resolver").(Resolver)
		balancer := NewRoundRobinBalancer()

		return NewManager(connList, resolver, balancer)
	})
}
