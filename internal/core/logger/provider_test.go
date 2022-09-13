package logger

import (
	"os"
	"testing"
	"time"

	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type LoggerProviderTestSuite struct {
	suite.Suite
	container core.Container
	provider  *Provider
}

func (s *LoggerProviderTestSuite) SetupTest() {
	s.container = core.NewApp()
	s.provider = new(Provider)
}

func (s *LoggerProviderTestSuite) TestRegisterLevel() {
	s.provider.Register(s.container)

	level, err := s.container.Get("logger.level")
	s.Require().NoError(err)
	s.Require().Equal(defaultLogLevel, level)
}

func (s *LoggerProviderTestSuite) TestRegisterCustomLevel() {
	os.Setenv(levelEnvName, "debug")
	defer os.Unsetenv(levelEnvName)

	s.provider.Register(s.container)

	level, err := s.container.Get("logger.level")
	s.Require().NoError(err)
	s.Require().Equal(logrus.DebugLevel, level)
}

func (s *LoggerProviderTestSuite) TestRegisterIncorrectCustomTimeout() {
	incorrectLevel := "debug_test"

	os.Setenv(levelEnvName, incorrectLevel)
	defer os.Unsetenv(levelEnvName)

	s.provider.Register(s.container)

	level, err := s.container.Get("logger.level")
	s.Require().NoError(err)
	s.Require().Equal(defaultLogLevel, level)
}

func (s *LoggerProviderTestSuite) TestRegisterFormat() {
	s.provider.Register(s.container)

	formatter, err := s.container.Get("logger.format")
	s.Require().NoError(err)
	s.Require().Equal(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano}, formatter)
}

func (s *LoggerProviderTestSuite) TestRegisterJsonFormat() {
	os.Setenv(formatEnvName, "json")
	defer os.Unsetenv(formatEnvName)

	s.provider.Register(s.container)

	formatter, err := s.container.Get("logger.format")
	s.Require().NoError(err)
	s.Require().Equal(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano}, formatter)
}

func (s *LoggerProviderTestSuite) TestRegisterIncorrectCustomFormat() {
	incorrectFormat := "xml"

	os.Setenv(formatEnvName, incorrectFormat)
	defer os.Unsetenv(formatEnvName)

	s.provider.Register(s.container)

	formatter, err := s.container.Get("logger.format")
	s.Require().NoError(err)
	s.Require().Equal(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano}, formatter)
}

func (s *LoggerProviderTestSuite) TestLoggerService() {
	s.provider.Register(s.container)

	logger, err := s.container.Get("logger")

	s.Require().NoError(err, "It should provide logger service")
	s.Require().Implements((*logrus.FieldLogger)(nil), logger, "Logger should implement Logger interface")
}

func TestLoggerProviderTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerProviderTestSuite))
}
