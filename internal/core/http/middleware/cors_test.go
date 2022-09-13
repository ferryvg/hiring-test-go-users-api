package middleware

import (
	"testing"
	"time"

	"github.com/ferryvg/hiring-test-go-users-api/internal/core/http"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type CorsTestSuite struct {
	suite.Suite
	middleware http.Middleware
}

func (s *CorsTestSuite) SetupTest() {
	s.middleware = CorsMiddleware
}

func (s *CorsTestSuite) TestGET() {
	ctx := new(fasthttp.RequestCtx)

	origin := "http://another.com/"
	ctx.Request.URI().SetHost("test.com")
	ctx.Request.SetRequestURI("http://test.com/test/path")
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Request.Header.Add(fasthttp.HeaderOrigin, origin)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	called := false
	next := func(ctx *fasthttp.RequestCtx) {
		time.Sleep(10 * time.Millisecond)
		called = true
	}

	handler := s.middleware(next)
	handler(ctx)

	s.Require().True(called)
	s.Require().Equal(string(ctx.Response.Header.Peek(fasthttp.HeaderAccessControlAllowOrigin)), origin)
	s.Require().Equal(string(ctx.Response.Header.Peek(fasthttp.HeaderAccessControlAllowCredentials)), "true")
}

func (s *CorsTestSuite) TestGETWithoutOrigin() {
	ctx := new(fasthttp.RequestCtx)

	ctx.Request.URI().SetHost("test.com")
	ctx.Request.SetRequestURI("http://test.com/test/path")
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	called := false
	next := func(ctx *fasthttp.RequestCtx) {
		time.Sleep(10 * time.Millisecond)
		called = true
	}

	handler := s.middleware(next)
	handler(ctx)

	s.Require().True(called)
	s.Require().Empty(string(ctx.Response.Header.Peek(fasthttp.HeaderAccessControlAllowOrigin)))
	s.Require().Empty(string(ctx.Response.Header.Peek(fasthttp.HeaderAccessControlAllowCredentials)))
}

func (s *CorsTestSuite) TestOPTIONS() {
	ctx := new(fasthttp.RequestCtx)

	origin := "http://another.com/"
	ctx.Request.URI().SetHost("test.com")
	ctx.Request.SetRequestURI("http://test.com/test/path")
	ctx.Request.Header.SetMethod(fasthttp.MethodOptions)
	ctx.Request.Header.Add(fasthttp.HeaderOrigin, origin)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	called := false
	next := func(ctx *fasthttp.RequestCtx) {
		time.Sleep(10 * time.Millisecond)
		called = true
	}

	handler := s.middleware(next)
	handler(ctx)

	s.Require().False(called)
	s.Require().Equal(string(ctx.Response.Header.Peek(fasthttp.HeaderAccessControlAllowOrigin)), origin)
	s.Require().Equal(string(ctx.Response.Header.Peek(fasthttp.HeaderAccessControlAllowCredentials)), "true")
}

func (s *CorsTestSuite) TestOPTIONSWithoutOrigin() {
	ctx := new(fasthttp.RequestCtx)

	ctx.Request.URI().SetHost("test.com")
	ctx.Request.SetRequestURI("http://test.com/test/path")
	ctx.Request.Header.SetMethod(fasthttp.MethodOptions)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	called := false
	next := func(ctx *fasthttp.RequestCtx) {
		time.Sleep(10 * time.Millisecond)
		called = true
	}

	handler := s.middleware(next)
	handler(ctx)

	s.Require().False(called)
	s.Require().Empty(string(ctx.Response.Header.Peek(fasthttp.HeaderAccessControlAllowOrigin)))
	s.Require().Empty(string(ctx.Response.Header.Peek(fasthttp.HeaderAccessControlAllowCredentials)))
}

func TestCorsTestSuite(t *testing.T) {
	suite.Run(t, new(CorsTestSuite))
}
