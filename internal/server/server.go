package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adepte-myao/test_parser/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	http.Server
	config          *ServerConfig
	logger          *logrus.Logger
	router          *mux.Router
	taskPageHandler *handlers.TaskPageHandler
	linksHandler    *handlers.LinksHandler
	solutionHandler *handlers.SolutionHandler
}

func NewServer(config *ServerConfig) *Server {
	serv := &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}

	serv.taskPageHandler = handlers.NewTaskPageHandler(serv.logger)
	serv.linksHandler = handlers.NewLinksHandler(serv.logger, config.BaseLink)
	serv.solutionHandler = handlers.NewSolutionHandler(serv.logger, config.BaseLink)

	return serv
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()
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
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/site", s.taskPageHandler.Handle)
	s.router.HandleFunc("/links", s.linksHandler.Handle)
	s.router.HandleFunc("/solution", s.solutionHandler.Handle)
}

func (s *Server) congfigureServer() {
	s.Addr = s.config.BindAddr
	s.Handler = s.router
	s.IdleTimeout = 120 * time.Second
	s.ReadTimeout = 3 * time.Second
	s.WriteTimeout = 0 * time.Second
}
