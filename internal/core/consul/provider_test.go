package consul

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ProviderTestSuite struct {
	suite.Suite
	container core.Container
	provider  *Provider
}

func (s *ProviderTestSuite) SetupTest() {
	s.container = core.NewApp()
	s.container.Set("saas.datacenter", "dc1")
	s.provider = new(Provider)
}

func (s *ProviderTestSuite) TestRegister() {
	s.provider.Register(s.container)

	client, err := s.container.Get("consul.client")
	s.Require().NoError(err)
	s.Require().IsType((*api.Client)(nil), client)
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
