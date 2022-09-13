package consul

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/hashicorp/consul/api"
)

type Provider struct{}

func (p *Provider) Register(c core.Container) {
	c.Set("consul.client", func(c core.Container) interface{} {
		dc := c.MustGet("saas.datacenter").(string)

		config := api.DefaultConfig()
		config.Datacenter = dc

		client, err := api.NewClient(config)
		if err != nil {
			panic(err)
		}

		return client
	})
}
