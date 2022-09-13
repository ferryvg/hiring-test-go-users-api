package metrics

import (
	"testing"

	"github.com/buaazp/fasthttprouter"
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

	s.container.Set("http.router", func(c core.Container) interface{} {
		return fasthttprouter.New()
	})

	s.provider = new(Provider)
}

func (s *ProviderTestSuite) TestRegisterRoute() {
	s.provider.Register(s.container)

	routes := s.container.MustGet("http.router").(*fasthttprouter.Router)

	handler, _ := routes.Lookup("GET", "/metrics", nil)

	s.Require().NotNil(handler)
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
