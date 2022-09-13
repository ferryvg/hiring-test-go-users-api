package config_test

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type BuilderTestSuite struct {
	suite.Suite
	logger     logrus.FieldLogger
	loggerHook *test.Hook
	builder    config.Builder
}

func (s *BuilderTestSuite) SetupTest() {
	s.logger, s.loggerHook = test.NewNullLogger()

	s.builder = config.NewBuilder(s.logger)
}

func (s *BuilderTestSuite) TestBuild() {
	expConf := s.getDefaultConf()

	conf, err := s.builder.Build("")

	s.Require().NoError(err)
	s.Require().Equal(expConf, conf)
}

func (s *BuilderTestSuite) TestBuildFromEnvVars() {
	expConf := s.getTestConf()

	s.setupEnvVars(expConf)
	defer s.clearEnvVars()

	conf, err := s.builder.Build("")

	s.Require().NoError(err)
	s.Require().Equal(expConf, conf)
}

func (s *BuilderTestSuite) TestBuildFromFile() {
	expConf := s.getTestConf()

	conf, err := s.builder.Build("./../../test/test_config.yaml")

	s.Require().NoError(err)
	s.Require().Equal(expConf, conf)
}

func (s *BuilderTestSuite) setupEnvVars(conf *config.AppConfig) {
	_ = os.Setenv("APP_MYSQL_DATABASE", conf.Mysql.Database)
	_ = os.Setenv("APP_MYSQL_USERNAME", conf.Mysql.Username)
	_ = os.Setenv("APP_MYSQL_PASSWORD", conf.Mysql.Password)
	_ = os.Setenv("APP_JWT_TTL", conf.Jwt.TTL.String())
}

func (s *BuilderTestSuite) clearEnvVars() {
	vars := []string{
		"APP_MYSQL_DATABASE",
		"APP_MYSQL_USERNAEM",
		"APP_MYSQL_PASSWORD",
		"APP_JWT_TTL",
	}

	for _, name := range vars {
		_ = os.Unsetenv(name)
	}
}

func (s *BuilderTestSuite) getDefaultConf() *config.AppConfig {
	return &config.AppConfig{
		Mysql: config.MysqlConfig{
			Database: "users_api",
			Username: "root",
			Password: "scout",
		},
		Jwt: config.JwtConfig{
			TTL: time.Hour * 720,
		},
	}
}

func (s *BuilderTestSuite) getTestConf() *config.AppConfig {
	return &config.AppConfig{
		Mysql: config.MysqlConfig{
			Database: "some_db_name",
			Username: "db_user",
			Password: "123456",
		},
		Jwt: config.JwtConfig{
			TTL: time.Hour * 24,
		},
	}
}

func TestConfigBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
