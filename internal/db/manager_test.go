package db_test

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/config"
	"github.com/ferryvg/hiring-test-go-users-api/internal/db"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type mockResolver struct {
	mock.Mock
}

func (m *mockResolver) Init() error {
	return m.Called().Error(0)
}

func (m *mockResolver) Shutdown() {
	m.Called()
}

type ManagerTestSuite struct {
	suite.Suite
	resolver         *mockResolver
	logger           logrus.FieldLogger
	loggerHook       *test.Hook
	connections      db.ConnList
	clusterNodeStore config.ClusterNodeStore
	manager          db.Manager
	balancer         db.Balancer
}

func (s *ManagerTestSuite) SetupTest() {
	s.logger, s.loggerHook = test.NewNullLogger()
	s.resolver = new(mockResolver)
	s.balancer = db.NewRoundRobinBalancer()
	s.clusterNodeStore = config.NewClusterNodeStore()

	factory := db.NewMysqlConnFactory("users_api", "root", "scout")
	s.connections = db.NewConnList(factory, s.clusterNodeStore, s.logger)

	s.manager = db.NewManager(s.connections, s.resolver, s.balancer)
}

func (s *ManagerTestSuite) TestShutdown() {
	nodes := []string{"db"}

	s.connections.SetNodes(nodes)

	s.resolver.On("Shutdown").Once()
	s.manager.Shutdown()

	s.Require().Empty(s.connections.Connections())
}

func (s *ManagerTestSuite) TestGetDb() {
	nodes := []string{"db"}

	s.connections.SetNodes(nodes)

	next, err := s.manager.GetDB()
	conn, exists := s.connections.Connection(nodes[0])
	s.Require().True(exists)
	s.Require().Equal(conn, next)
	s.Require().NoError(err)
}

func TestManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}
