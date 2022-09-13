package http

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
)

type testHandler struct {
	username string
}

func (h *testHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello, %s", h.username)
}

func newTestHandler(name string) fasthttp.RequestHandler {
	handler := fasthttpadaptor.NewFastHTTPHandler(&testHandler{username: name})
	return handler
}

type ProviderTestSuite struct {
	suite.Suite
	container core.Container
	provider  *Provider

	logger     *logrus.Logger
	loggerHook *test.Hook
}

func (s *ProviderTestSuite) SetupTest() {
	s.container = core.NewApp()
	s.provider = NewProvider(":0")

	s.logger, s.loggerHook = test.NewNullLogger()

	s.container.Set("logger", func(c core.Container) interface{} {
		return s.logger
	})
}

func (s *ProviderTestSuite) TestRegisterAddr() {
	s.provider.Register(s.container)

	addr, err := s.container.Get("http.addr")
	s.Require().NoError(err)
	s.Require().Equal(s.provider.addr, addr)
}

func (s *ProviderTestSuite) TestRegisterCustomAddr() {
	expAddr := "0.0.0.0:8081"
	os.Setenv(serverAddrEnvName, expAddr)
	defer os.Unsetenv(serverAddrEnvName)

	s.provider.Register(s.container)

	addr, err := s.container.Get("http.addr")
	s.Require().NoError(err)
	s.Require().Equal(expAddr, addr)
}

func (s *ProviderTestSuite) TestRegisterRoutes() {
	s.provider.Register(s.container)

	value, err := s.container.Get("http.router")
	s.Require().NoError(err)

	routes := value.(*fasthttprouter.Router)
	s.Require().Equal(routes, fasthttprouter.New())
}

func (s *ProviderTestSuite) TestRegisterServer() {
	s.provider.Register(s.container)

	server, err := s.container.Get("http.server")
	s.Require().NoError(err)
	s.Require().IsType((*fasthttp.Server)(nil), server)
}

func (s *ProviderTestSuite) TestRegisterListener() {
	s.container.Set("http.addr", s.getServerAddr())
	s.provider.Register(s.container)

	listener, err := s.container.Get("http.listener")
	s.Require().NoError(err)
	s.Require().Implements((*net.Listener)(nil), listener)
}

func (s *ProviderTestSuite) TestBoot() {
	s.provider.Register(s.container)

	s.container.MustExtend("http.router", func(old interface{}, c core.Container) interface{} {
		routes := old.(*fasthttprouter.Router)
		routes.GET("/test", newTestHandler("User"))

		return routes
	})

	addr := s.getServerAddr()
	s.container.Set("http.addr", addr)

	err := s.provider.Boot(s.container)
	s.Require().NoError(err)

	url := fmt.Sprintf("http://%s/test", addr)
	res, err := http.Get(url)
	s.Require().NoError(err)
	defer res.Body.Close()

	s.Require().Equal(http.StatusOK, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	s.Require().NoError(err)

	s.Require().Equal("Hello, User", string(body))

	s.provider.Shutdown(s.container)
}

func (s *ProviderTestSuite) getServerAddr() string {
	listener, _ := net.Listen("tcp", ":0")
	addr := listener.Addr().String()
	listener.Close()

	return addr
}

func TestServerProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
