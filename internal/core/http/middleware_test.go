package http

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type MiddlewareTestSuite struct {
	suite.Suite
	firstMiddleware        Middleware
	firstMiddlewareCalled  bool
	secondMiddleware       Middleware
	secondMiddlewareCalled bool
	request                *fasthttp.Request
	ctx                    *fasthttp.RequestCtx
}

func (s *MiddlewareTestSuite) SetupTest() {
	s.firstMiddlewareCalled = false
	s.firstMiddleware = func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			s.firstMiddlewareCalled = true
			next(ctx)
		}
	}

	s.secondMiddlewareCalled = false
	s.secondMiddleware = func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			s.secondMiddlewareCalled = true
			next(ctx)
		}
	}

	s.request = fasthttp.AcquireRequest()
	s.ctx = new(fasthttp.RequestCtx)
	s.ctx.Init(s.request, nil, nil)
}

func (s *MiddlewareTestSuite) TearDownTest() {
	fasthttp.ReleaseRequest(s.request)
}

func (s *MiddlewareTestSuite) TestBuildHandler() {
	handlerCalled := false
	handler := func(_ *fasthttp.RequestCtx) {
		handlerCalled = true
	}

	endpoint := BuildHandler(handler, s.firstMiddleware, s.secondMiddleware)
	endpoint(s.ctx)

	s.Require().True(s.secondMiddlewareCalled)
	s.Require().True(s.firstMiddlewareCalled)
	s.Require().True(handlerCalled)
}

func TestEndpointBuilderSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
