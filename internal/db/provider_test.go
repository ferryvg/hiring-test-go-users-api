package db_test

import (
	"context"
	"github.com/ferryvg/hiring-test-go-users-api/internal/config"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core/saas"
	"github.com/ferryvg/hiring-test-go-users-api/internal/db"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MockRegistry struct {
	mock.Mock
}

func (r *MockRegistry) Get(ctx context.Context, service string, tags []string, waitIdx uint64) ([]string, uint64, error) {
	args := r.Called(ctx, service, tags, waitIdx)

	return args.Get(0).([]string), args.Get(1).(uint64), args.Error(2)
}

type ProviderSqlTestSuite struct {
	suite.Suite
	container  core.Container
	logger     logrus.FieldLogger
	loggerHook *test.Hook
	provider   *db.Provider
	saas       *saas.Provider
}

func (s *ProviderSqlTestSuite) SetupTest() {
	s.container = core.NewApp()
	s.saas.Register(s.container)

	s.logger, s.loggerHook = test.NewNullLogger()

	s.container.Set("logger", s.logger)
	s.container.Set("svc.config", new(config.AppConfig))

	factory := db.NewMysqlConnFactory("test", "test", "test")
	s.container.Set("svc.db.conn_factory", factory)

	connectionList := db.NewConnList(factory, config.NewClusterNodeStore(), s.logger)
	s.container.Set("svc.db.conn_list", connectionList)

	resolver := new(mockResolver)
	s.container.Set("sd.registry", new(MockRegistry))
	s.container.Set("svc.db.resolver", resolver)

	manager := new(db.Manager)
	s.container.Set("svc.db.manager", manager)

	s.provider = new(db.Provider)
}

func (s *ProviderSqlTestSuite) TestRegisterFactory() {
	s.provider.Register(s.container)

	factory, err := s.container.Get("svc.db.conn_factory")

	s.Require().NoError(err)
	s.Require().Implements((*db.ConnFactory)(nil), factory)
}

func (s *ProviderSqlTestSuite) TestRegisterConnList() {
	s.provider.Register(s.container)

	factory, err := s.container.Get("svc.db.conn_list")

	s.Require().NoError(err)
	s.Require().Implements((*db.ConnList)(nil), factory)
}

func (s *ProviderSqlTestSuite) TestRegisterResolver() {
	s.provider.Register(s.container)

	factory, err := s.container.Get("svc.db.resolver")

	s.Require().NoError(err)
	s.Require().Implements((*db.Resolver)(nil), factory)
}

func (s *ProviderSqlTestSuite) TestRegisterManager() {
	s.provider.Register(s.container)

	factory, err := s.container.Get("svc.db.manager")

	s.Require().NoError(err)
	s.Require().Implements((*db.Manager)(nil), factory)
}

func TestProviderSqlTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderSqlTestSuite))
}
