package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adepte-myao/test_parser/internal/config"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	http.Server
	config *config.Config
	logger *logrus.Logger
	router *mux.Router
}

func NewServer(config *config.Config, logger *logrus.Logger, router *mux.Router) *Server {
	serv := &Server{
		config: config,
		logger: logger,
		router: router,
	}

	return serv
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.congfigureServer()

	errChan := make(chan error, 1)
	go func() {
		s.logger.Info("Server started")

		err := s.ListenAndServe()
		if err != nil {
			errChan <- err
			return
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	select {
	case sig := <-sigChan:
		s.logger.Info("Received terminate, graceful shutdown. Signal:", sig)
		tc, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancelFunc()

		s.Shutdown(tc)
	case err := <-errChan:
		return err
	}

	return nil
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.Server.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

func (s *Server) RegisterHandler(path string, handler http.HandlerFunc) {
	s.router.HandleFunc(path, handler)
}

func (s *Server) congfigureServer() {
	s.Addr = s.config.Server.BindAddr
	s.Handler = s.router
	s.IdleTimeout = 120 * time.Second
	s.ReadTimeout = 3 * time.Second
	s.WriteTimeout = 0 * time.Second
}

func (s *Server) Ping(rw http.ResponseWriter, r *http.Request) {
	s.logger.Info("Ping request received")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hello from proxy!"))
}
