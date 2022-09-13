package config_test

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/config"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ProviderTestSuite struct {
	suite.Suite
	container core.Container
	provider  *config.Provider
}

func (s *ProviderTestSuite) SetupTest() {
	s.container = core.NewApp()

	s.container.Set("logger", func(c core.Container) interface{} {
		logger, _ := test.NewNullLogger()

		return logger
	})
	s.container.Set("svc.config_file", func(c core.Container) interface{} {
		var confFile = ""

		return &confFile
	})

	s.provider = new(config.Provider)
}

func (s *ProviderTestSuite) TestRegisterBuilder() {
	s.provider.Register(s.container)

	builder, err := s.container.Get("svc.config.builder")

	s.Require().NoError(err)
	s.Require().Implements((*config.Builder)(nil), builder)
}

func (s *ProviderTestSuite) TestRegisterConfig() {
	s.provider.Register(s.container)

	conf, err := s.container.Get("svc.config")

	s.Require().NoError(err)

	_, ok := conf.(*config.AppConfig)
	s.Require().True(ok)
}

func TestConfigProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
