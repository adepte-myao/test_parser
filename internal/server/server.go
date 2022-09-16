package server

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adepte-myao/test_parser/internal/html"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	http.Server
	config     *ServerConfig
	logger     *logrus.Logger
	router     *mux.Router
	htmlParser *html.Parser
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		config:     config,
		logger:     logrus.New(),
		router:     mux.NewRouter(),
		htmlParser: html.NewParser(),
	}
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
		tc, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
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
	s.router.HandleFunc("/", s.handleHmtl)
	s.router.HandleFunc("/file", s.handleFile)
}

func (s *Server) congfigureServer() {
	s.Addr = s.config.BindAddr
	s.Handler = s.router
	s.IdleTimeout = 120 * time.Second
	s.ReadTimeout = 3 * time.Second
	s.WriteTimeout = 3 * time.Second
}

func (s *Server) handleHmtl(rw http.ResponseWriter, r *http.Request) {
	s.logger.Info("Received handle simple HTML request")

	resp, err := http.Get("https://tests24.ru/?iter=3&test=726")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte("Response from given source is not OK"))
		return
	}

	_, err = io.Copy(rw, resp.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (s *Server) handleFile(rw http.ResponseWriter, r *http.Request) {
	s.logger.Info("Received handleFile request")

	data, err := os.ReadFile("src.html")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	dataStr := string(data)

	tasks := s.htmlParser.ParseHtml(dataStr)

	rw.WriteHeader(http.StatusOK)
	for _, task := range tasks {
		rw.Write([]byte(task.Question))

		rw.Write([]byte("\n"))
		for _, option := range task.Options {
			rw.Write([]byte("\t"))
			rw.Write([]byte(option))
			rw.Write([]byte("\n"))
		}

		rw.Write([]byte("\tRight answer: "))
		rw.Write([]byte(task.Answer))

		rw.Write([]byte("\n\n"))
	}
}
