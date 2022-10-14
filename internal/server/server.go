package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adepte-myao/test_parser/internal/config"
	"github.com/adepte-myao/test_parser/internal/handlers"
	"github.com/adepte-myao/test_parser/internal/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	http.Server
	config          *config.Config
	logger          *logrus.Logger
	router          *mux.Router
	linksHandler    *handlers.LinksHandler
	solutionHandler *handlers.SolutionHandler
	store           *storage.Store
}

func NewServer(config *config.Config) *Server {
	serv := &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}

	return serv
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	if err := s.configureStore(); err != nil {
		return err
	}
	defer s.store.Close()

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

		defer s.store.Close()
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

func (s *Server) configureRouter() {
	s.linksHandler = handlers.NewLinksHandler(s.logger, s.config.Server.BaseLink, s.store)
	s.solutionHandler = handlers.NewSolutionHandler(s.logger, s.config.Server.BaseLink, s.store)

	s.router.HandleFunc("/links", s.linksHandler.Handle)
	s.router.HandleFunc("/solution", s.solutionHandler.Handle)
	s.router.HandleFunc("/ping", s.ping)
}

func (s *Server) congfigureServer() {
	s.Addr = s.config.Server.BindAddr
	s.Handler = s.router
	s.IdleTimeout = 120 * time.Second
	s.ReadTimeout = 3 * time.Second
	s.WriteTimeout = 0 * time.Second
}

func (s *Server) configureStore() error {
	st := storage.NewStore((*storage.StoreConfig)(&s.config.Store), s.logger)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}

func (s *Server) ping(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hello from proxy!"))
}
