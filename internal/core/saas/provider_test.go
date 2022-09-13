package saas

import (
	"os"
	"testing"

	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/stretchr/testify/suite"
)

type ProviderTestSuite struct {
	suite.Suite
	container core.Container
	provider  *Provider
}

func (s *ProviderTestSuite) SetupTest() {
	s.container = core.NewApp()
	s.provider = new(Provider)
}

func (s *ProviderTestSuite) TestRegisterDatacenter() {
	s.provider.Register(s.container)

	dc, err := s.container.Get("saas.datacenter")
	s.Require().NoError(err)
	s.Require().Equal(defaultDatacenter, dc)
}

func (s *ProviderTestSuite) TestRegisterCustomDatacenter() {
	expDc := "custom1"

	os.Setenv(datacenterEnvName, expDc)
	defer os.Unsetenv(datacenterEnvName)

	s.provider.Register(s.container)

	dc, err := s.container.Get("saas.datacenter")
	s.Require().NoError(err)
	s.Require().Equal(expDc, dc)
}

func (s *ProviderTestSuite) TestRegisterCluster() {
	s.provider.Register(s.container)

	cluster, err := s.container.Get("saas.cluster")
	s.Require().NoError(err)
	s.Require().Equal(defaultCluster, cluster)
}

func (s *ProviderTestSuite) TestRegisterCustomCluster() {
	expCluster := "test"
	os.Setenv(clusterEnvName, expCluster)
	defer os.Unsetenv(clusterEnvName)

	s.provider.Register(s.container)

	cluster, err := s.container.Get("saas.cluster")
	s.Require().NoError(err)
	s.Require().Equal(expCluster, cluster)
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
