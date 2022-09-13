package dal_test

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"strings"
	"testing"
	"time"

	"github.com/ferryvg/hiring-test-go-users-api/internal/config"
	"github.com/ferryvg/hiring-test-go-users-api/internal/dal"
	"github.com/ferryvg/hiring-test-go-users-api/internal/db"
	"github.com/ferryvg/hiring-test-go-users-api/internal/domain"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UsersManagerTestSuite struct {
	suite.Suite

	config    *config.JwtConfig
	dbManager db.Manager
	manager   dal.UsersManager

	logger     *logrus.Logger
	loggerHook *test.Hook

	ctx context.Context
}

func (s *UsersManagerTestSuite) SetupTest() {
	s.config = &config.JwtConfig{
		TTL: 2 * time.Second,
	}
	s.logger, s.loggerHook = test.NewNullLogger()

	s.dbManager = NewDbManager()

	s.manager = dal.NewUsersManager(s.config, s.dbManager, s.logger)

	s.ctx = context.Background()
}

func (s *UsersManagerTestSuite) TearDownTest() {
	conn, err := s.dbManager.GetDB()
	s.Require().NoError(err)

	_, err = conn.Exec("TRUNCATE TABLE user_roles")
	s.Require().NoError(err)

	_, err = conn.Exec("TRUNCATE TABLE users")
	s.Require().NoError(err)
}

func (s *UsersManagerTestSuite) TestCreateAndChangeRoles() {
	expUser := &domain.User{
		ID:     "test-create-user",
		Secret: "123",
		Roles: map[domain.Role]bool{
			domain.BasicRole: true,
			domain.AdminRole: false,
		},
	}

	err := s.manager.Create(s.ctx, expUser)
	s.Require().NoError(err)

	actUser, err := s.manager.Get(s.ctx, expUser.ID)
	s.Require().NoError(err)
	s.Require().NotNil(actUser)

	s.Require().Equal(expUser.ID, actUser.ID)

	err = bcrypt.CompareHashAndPassword([]byte(actUser.Secret), []byte(expUser.Secret))
	s.Require().NoError(err)

	expRoles := map[domain.Role]bool{
		domain.BasicRole: true,
	}

	s.Require().Equal(expRoles, actUser.Roles)

	newRoles := map[domain.Role]bool{
		domain.BasicRole: true,
		domain.AdminRole: true,
	}

	err = s.manager.ChangeRoles(s.ctx, expUser.ID, newRoles)
	s.Require().NoError(err)

	actUser, err = s.manager.Get(s.ctx, expUser.ID)
	s.Require().NoError(err)
	s.Require().NotNil(actUser)
	s.Require().Equal(newRoles, actUser.Roles)
}

func (s *UsersManagerTestSuite) TestCreateInvalid() {
	s.Require().Error(s.manager.Create(s.ctx, &domain.User{}))

	s.Require().Error(s.manager.Create(s.ctx, &domain.User{
		ID: strings.Repeat("s-", 500),
	}))

	s.Require().NoError(s.manager.Create(s.ctx, &domain.User{
		ID:     "test-invalid",
		Secret: strings.Repeat("123", 1000),
	}))

	s.Require().Error(s.manager.Create(s.ctx, &domain.User{
		ID:     "test-invalid",
		Secret: "already exists",
	}))
}

func (s *UsersManagerTestSuite) TestAuthenticate() {
	expUser := &domain.User{
		ID:     "test-auth-user",
		Secret: "55555",
		Roles: map[domain.Role]bool{
			domain.BasicRole: true,
		},
	}

	err := s.manager.Create(s.ctx, expUser)
	s.Require().NoError(err)

	token, err := s.manager.Authenticate(s.ctx, expUser.ID, expUser.Secret)
	s.Require().NoError(err)
	s.Require().NotEmpty(token)

	time.Sleep(time.Millisecond * 1500)

	tokenUser, err := s.manager.VerifyToken(s.ctx, token.AccessToken)
	s.Require().NoError(err)
	s.Require().Equal(expUser.ID, tokenUser.ID)

	time.Sleep(s.config.TTL)
	_, err = s.manager.VerifyToken(s.ctx, token.AccessToken)
	s.Require().Error(err)
}

func (s *UsersManagerTestSuite) TestAuthenticateFailure() {
	_, err := s.manager.Authenticate(s.ctx, "test-auth-failure", "invalid")
	s.Require().Error(err)
}

func (s *UsersManagerTestSuite) TestVerifyTokenFailure() {
	_, err := s.manager.VerifyToken(s.ctx, "invalid")
	s.Require().Error(err)
}

func (s *UsersManagerTestSuite) TestGetList() {
	expUsers := []domain.User{
		{
			ID:     "test-list-user-1",
			Secret: "123",
			Roles: map[domain.Role]bool{
				domain.BasicRole: true,
			},
		},
		{
			ID:     "test-list-user-2",
			Secret: "456",
			Roles: map[domain.Role]bool{
				domain.AdminRole: true,
			},
		},
	}

	for idx, user := range expUsers {
		s.Require().NoError(s.manager.Create(s.ctx, &user))

		expUsers[idx].Secret = ""
	}

	users, err := s.manager.GetList(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(expUsers, users)
}

func TestUsersManagerTestSuite(t *testing.T) {
	suite.Run(t, new(UsersManagerTestSuite))
}
