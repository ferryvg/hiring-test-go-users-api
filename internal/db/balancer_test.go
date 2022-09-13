package db_test

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BalancerTestSuite struct {
	suite.Suite
	balancer db.Balancer
}

func (s *BalancerTestSuite) SetupTest() {
	s.balancer = db.NewRoundRobinBalancer()
}

func (s *BalancerTestSuite) TestNext() {
	sqlDbs := []*sqlx.DB{new(sqlx.DB), new(sqlx.DB), new(sqlx.DB)}
	_, err := s.balancer.Next(sqlDbs)

	s.Require().NoError(err)
}

func (s *BalancerTestSuite) TestNextErr() {
	var sqlDbs []*sqlx.DB
	_, err := s.balancer.Next(sqlDbs)

	s.Require().Error(err)
}

func TestBalancerTestSuite(t *testing.T) {
	suite.Run(t, new(BalancerTestSuite))
}
