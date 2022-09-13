package core_test

import (
	"errors"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
	"time"
)

type TestProvider struct {
	mock.Mock
}

func (p *TestProvider) Register(container core.Container) {
	p.Called(container)
}

func (p *TestProvider) Boot(container core.Container) error {
	args := p.Called(container)

	return args.Error(0)
}

func (p *TestProvider) Shutdown(container core.Container) {
	p.Called(container)
}

func (p *TestProvider) Reconfigure(container core.Container) {
	p.Called(container)
}

type AppTestSuite struct {
	suite.Suite
	app *core.App
}

func (s *AppTestSuite) SetupTest() {
	s.app = core.NewApp()
}

func (s *AppTestSuite) TestRegister() {
	provider := new(TestProvider)
	provider.On("Register", s.app).Once()

	s.app.Register(provider)

	provider.AssertExpectations(s.T())
}

func (s *AppTestSuite) TestBoot() {
	provider := new(TestProvider)
	provider.On("Register", s.app).Once()
	provider.On("Boot", s.app).Return(nil).Once()

	s.app.Register(provider)
	err := s.app.Boot()

	s.Require().NoError(err)
	provider.AssertExpectations(s.T())
}

func (s *AppTestSuite) TestBootError() {
	provider := new(TestProvider)
	provider.On("Register", s.app).Once()

	expErr := errors.New("test")
	provider.On("Boot", s.app).Return(expErr).Once()

	s.app.Register(provider)
	err := s.app.Boot()

	s.Require().Equal(expErr, err)
	provider.AssertExpectations(s.T())
}

func (s *AppTestSuite) TestShutdown() {
	provider := new(TestProvider)
	provider.On("Register", s.app).Once()
	provider.On("Shutdown", s.app).Once()

	s.app.Register(provider)
	s.app.Shutdown()

	provider.AssertExpectations(s.T())
}

func (s *AppTestSuite) TestReconfigure() {
	provider := new(TestProvider)
	provider.On("Register", s.app).Once()
	provider.On("Reconfigure", s.app).Once()

	s.app.Register(provider)

	err := s.app.Reconfigure()
	s.Require().NoError(err)

	provider.AssertExpectations(s.T())
}

func (s *AppTestSuite) TestSet() {
	s.app.Set("number", 10)
	s.app.Set("string", "test")
	s.app.Set("bool", true)

	s.app.Set("service", func(c core.Container) interface{} {
		return new(TestProvider)
	})

	num, err := s.app.Get("number")
	s.Require().NoError(err)
	s.Require().Equal(10, num)

	str, err := s.app.Get("string")
	s.Require().NoError(err)
	s.Require().Equal("test", str)

	boolean, err := s.app.Get("bool")
	s.Require().NoError(err)
	s.Require().Equal(true, boolean)

	svc1, err := s.app.Get("service")
	s.Require().NoError(err)
	s.Require().IsType((*TestProvider)(nil), svc1)

	svc2, err := s.app.Get("service")
	s.Require().NoError(err)
	s.Require().IsType((*TestProvider)(nil), svc2)

	s.Require().Exactly(svc1, svc2)

	_, err = s.app.Get("unknown")
	s.Require().EqualError(err, "identifier 'unknown' is not defined")
}

func (s *AppTestSuite) TestFactory() {
	s.app.Factory("factory", func(c core.Container) interface{} {
		rand.Seed(time.Now().UnixNano())

		return rand.Intn(56200)
	})

	val1 := s.app.MustGet("factory").(int)
	val2 := s.app.MustGet("factory").(int)

	s.Require().NotEqual(val1, val2)
}

func (s *AppTestSuite) TestProtect() {
	expVal := func(c core.Container) interface{} {
		return "test"
	}
	s.app.Protect("protected", expVal)

	val := s.app.MustGet("protected").(func(c core.Container) interface{})

	s.Require().Equal(expVal(s.app), val(s.app))
}

func (s *AppTestSuite) TestHas() {
	s.app.Set("test", "test")

	s.Require().True(s.app.Has("test"))
	s.Require().False(s.app.Has("unknown"))
}

func (s *AppTestSuite) TestMustGet() {
	s.app.Set("test", "test")

	s.Require().Equal("test", s.app.MustGet("test"))
	s.Require().Panics(func() {
		s.app.MustGet("unknown")
	})
}

func (s *AppTestSuite) TestExtend() {
	s.app.Set("extend", func(c core.Container) interface{} {
		return []string{"first"}
	})

	err := s.app.Extend("extend", func(old interface{}, c core.Container) interface{} {
		val := old.([]string)
		val = append(val, "second")

		return val
	})

	s.Require().NoError(err)
	s.Require().Equal([]string{"first", "second"}, s.app.MustGet("extend"))
}

func (s *AppTestSuite) TestExtendNotDefinedError() {
	err := s.app.Extend("unknown", func(old interface{}, c core.Container) interface{} {
		return old
	})

	s.Require().EqualError(err, "identifier 'unknown' is not defined")
}

func (s *AppTestSuite) TestExtendNotServiceError() {
	s.app.Set("test", "test")

	err := s.app.Extend("test", func(old interface{}, c core.Container) interface{} {
		return old
	})

	s.Require().EqualError(err, "idenfier 'test' does not contains service definition")
}

func (s *AppTestSuite) TestMustExtend() {
	s.app.Set("test", func(c core.Container) interface{} {
		return []string{"first"}
	})
	s.app.MustExtend("test", func(old interface{}, c core.Container) interface{} {
		val := old.([]string)
		val = append(val, "second")

		return val
	})

	s.Require().Equal([]string{"first", "second"}, s.app.MustGet("test"))
}

func (s *AppTestSuite) TestMustExtendError() {
	s.Require().Panics(func() {
		s.app.MustExtend("unknown", func(old interface{}, c core.Container) interface{} {
			return old
		})
	})
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}
