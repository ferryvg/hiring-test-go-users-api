package middleware

import (
	"testing"
	"time"

	"github.com/ferryvg/hiring-test-go-users-api/internal/core/http"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type LoggerTestSuite struct {
	suite.Suite
	logger     logrus.FieldLogger
	loggerHook *test.Hook
	middleware http.Middleware
	ctx        *fasthttp.RequestCtx
}

func (s *LoggerTestSuite) SetupTest() {
	s.logger, s.loggerHook = test.NewNullLogger()
	s.middleware = Logger(s.logger)

	s.ctx = new(fasthttp.RequestCtx)
}

func (s *LoggerTestSuite) TearDownTest() {
	fasthttp.ReleaseRequest(&s.ctx.Request)
	fasthttp.ReleaseResponse(&s.ctx.Response)
}

func (s *LoggerTestSuite) TestWork() {
	s.ctx.Request.URI().SetHost("test.com")
	s.ctx.Request.SetRequestURI("http://test.com/test/path")
	s.ctx.Request.Header.SetMethod("GET")
	s.ctx.Response.SetStatusCode(fasthttp.StatusOK)

	called := false
	next := func(ctx *fasthttp.RequestCtx) {
		time.Sleep(10 * time.Millisecond)
		called = true
	}

	handler := s.middleware(next)
	handler(s.ctx)

	s.Require().True(called)
	s.Require().NotNil(s.loggerHook.LastEntry())
	s.Require().Equal(logrus.InfoLevel, s.loggerHook.LastEntry().Level)
	s.Require().Equal("finished call", s.loggerHook.LastEntry().Message)
	s.Require().Equal(string(s.ctx.Request.URI().Path()), s.loggerHook.LastEntry().Data["http.endpoint"])
	s.Require().Equal(string(s.ctx.Request.Header.Method()), s.loggerHook.LastEntry().Data["http.method"])
	s.Require().Equal(s.ctx.Response.StatusCode(), s.loggerHook.LastEntry().Data["http.code"])
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
