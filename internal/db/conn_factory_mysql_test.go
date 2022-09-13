package db_test

import (
	"github.com/ferryvg/hiring-test-go-users-api/internal/db"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CFMTestSuite struct {
	suite.Suite
	cfm db.ConnFactory
}

func (s *CFMTestSuite) SetupTest() {
	s.cfm = db.NewMysqlConnFactory("users_api", "root", "scout")
}

func (s *CFMTestSuite) TestCreate() {
	sql, err := s.cfm.Create("db:3306")

	s.Require().NoError(err)
	s.NotNil(sql)
}

func TestCFMTestSuite(t *testing.T) {
	suite.Run(t, new(CFMTestSuite))
}
