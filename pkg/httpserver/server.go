package httpserver

import (
	"time"

	"github.com/valyala/fasthttp"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = ":5500"
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server          *fasthttp.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(handler fasthttp.RequestHandler, opts ...Option) *Server {
	httpServer := &fasthttp.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start(defaultAddr)

	return s
}

func (s *Server) start(addr string) {
	go func() {
		s.notify <- s.server.ListenAndServe(addr)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown()
}
