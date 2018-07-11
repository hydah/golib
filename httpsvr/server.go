package httpsvr

import (
	"net/http"
	"runtime"
	"time"

	"github.com/simplejia/utils"
	"github.com/tylerb/graceful"

	"github.com/hydah/golib/logger"
)

type HTTPServer struct {
	// Engine
	engine *Engine

	// Timeout is the duration to allow outstanding requests to survive
	// before forcefully terminating them.
	DelayTimeout time.Duration

	// maximum duration before timing out read of the request
	ReadTimeout time.Duration

	// maximum duration before timing out write of the response
	WriteTimeout time.Duration

	// enable hijact signal
	enableHijactSignal bool
}

func NewHTTPServer() *HTTPServer {
	s := &HTTPServer{
		engine: New(),

		DelayTimeout: 1 * time.Second,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	return s
}

func (s *HTTPServer) Init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (s *HTTPServer) EnableHijactSignal() {
	s.enableHijactSignal = true
}

func (s *HTTPServer) UseLogger() {
	s.engine.Use(Logger())
}

func (s *HTTPServer) UseRecovery() {
	s.engine.Use(Recovery())
}

func (s *HTTPServer) Get(pattern string, c IController) {
	s.engine.GETController(pattern, c)
}

func (s *HTTPServer) Post(pattern string, c IController) {
	s.engine.POSTController(pattern, c)
}

func (s *HTTPServer) Put(pattern string, c IController) {
	s.engine.PUTController(pattern, c)
}

func (s *HTTPServer) AddRoute(method string, pattern string, c IController) {
	switch method {
	case "GET":
		s.engine.GETController(pattern, c)
	case "POST":
		s.engine.POSTController(pattern, c)
	default:
		logger.Error("method [%s] mismatch", method)
	}
}

func (s *HTTPServer) Static(path, dir string) {
	s.engine.Static(path, dir)
}

func (s *HTTPServer) Run(hostport string) error {
	if s.enableHijactSignal {
		waitTime := time.Second * 5
		startTime := time.Second * 5
		err := utils.ListenAndServeWithTimeout(
			hostport,
			s.engine,
			waitTime,
			startTime,
		)
		if err != nil {
			logger.Error("%v, host: %v", err, utils.LocalIp)
			return err
		}

	} else {
		srv := &graceful.Server{
			Server: &http.Server{Addr: hostport, Handler: s.engine},
		}
		srv.Timeout = s.DelayTimeout
		srv.Server.ReadTimeout = s.ReadTimeout
		srv.Server.WriteTimeout = s.WriteTimeout
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("%v, host: %v", err, utils.LocalIp)
			return err
		}
	}

	return nil
}

func (s *HTTPServer) RunAsHttps(hostport, cert, key string) error {
	srv := &graceful.Server{
		Server: &http.Server{Addr: hostport, Handler: s.engine},
	}
	srv.Timeout = s.DelayTimeout
	srv.Server.ReadTimeout = s.ReadTimeout
	srv.Server.WriteTimeout = s.WriteTimeout
	if err := srv.ListenAndServeTLS(cert, key); err != nil {
		logger.Error("%v, host: %v", err, utils.LocalIp)
		return err
	}
	return nil
}
