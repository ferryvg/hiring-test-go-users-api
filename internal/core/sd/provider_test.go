package sd

import (
	"testing"

	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/suite"
)

type ProviderTestSuite struct {
	suite.Suite
	container core.Container
	provider  *Provider
}

func (s *ProviderTestSuite) SetupTest() {
	s.container = core.NewApp()

	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		s.T().Fatalf("Couldn't create consul client: %s", err)
	}

	s.container.Set("consul.client", consul)
	s.container.Set("saas.datacenter", "test")
	s.container.Set("saas.cluster", "test")
}

func (s *ProviderTestSuite) TestRegisterRegistry() {
	s.provider.Register(s.container)

	registry, err := s.container.Get("sd.registry")
	s.Require().NoError(err)
	s.Require().Implements((*Registry)(nil), registry)
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
