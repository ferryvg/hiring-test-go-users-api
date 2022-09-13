package saas

import (
	"os"

	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
)

const (
	datacenterEnvName = "SAAS_DC"
	clusterEnvName    = "SAAS_CLUSTER"
)

const (
	defaultDatacenter = "dc1"
	defaultCluster    = "dev"
)

type Provider struct{}

func (p *Provider) Register(c core.Container) {
	p.registerDatacenter(c)
	p.registerCluster(c)
}

func (p *Provider) registerDatacenter(c core.Container) {
	c.Set("saas.datacenter", func(c core.Container) interface{} {
		dc, exist := os.LookupEnv(datacenterEnvName)
		if exist && dc != "" {
			return dc
		}

		return defaultDatacenter
	})
}

func (p *Provider) registerCluster(c core.Container) {
	c.Set("saas.cluster", func(c core.Container) interface{} {
		cluster, exist := os.LookupEnv(clusterEnvName)
		if exist && cluster != "" {
			return cluster
		}

		return defaultCluster
	})
}
