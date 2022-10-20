package main

import (
	"fmt"
	"os"

	"github.com/adepte-myao/test_parser/internal/handlers"
	"github.com/adepte-myao/test_parser/internal/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/adepte-myao/test_parser/internal/config"
	"github.com/adepte-myao/test_parser/internal/server"
	"gopkg.in/yaml.v2"
)

func main() {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		fmt.Println("[ERROR]: ", err)
		return
	}
	defer f.Close()

	var cfg config.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("[ERROR]: ", err)
		return
	}

	logger := logrus.New()
	router := mux.NewRouter()
	store := storage.NewStore(&cfg.Store, logger)
	if err = store.Open(); err != nil {
		fmt.Println("[ERROR]: ", err)
		return
	}

	server := server.NewServer(&cfg, logger, router)

	server.RegisterHandler("/link", handlers.NewLinksHandler(logger, cfg.Server.BaseLink, store).Handle)
	server.RegisterHandler("/sitemap", handlers.NewSitemapHandler(logger, cfg.Server.BaseLink, store).Handle)
	server.RegisterHandler("/solution", handlers.NewSolutionHandler(logger, cfg.Server.BaseLink, store).Handle)
	server.RegisterHandler("/ping", server.Ping)

	err = server.Start()
	if err != nil {
		fmt.Println("[ERROR]: ", err)
		return
	}
}
