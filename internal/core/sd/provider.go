package sd

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/hashicorp/consul/api"
)

type Provider struct{}

func (p *Provider) Register(c core.Container) {
	p.registerRegistry(c)
}

func (p *Provider) registerRegistry(c core.Container) {
	c.Set("sd.registry", func(c core.Container) interface{} {
		consul := c.MustGet("consul.client").(*api.Client)
		dc := c.MustGet("saas.datacenter").(string)
		cluster := c.MustGet("saas.cluster").(string)

		return NewRegistry(consul, dc, cluster)
	})
}
