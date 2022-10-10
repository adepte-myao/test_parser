package main

import (
	"fmt"
	"os"

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

	server := server.NewServer(&cfg)
	err = server.Start()
	if err != nil {
		fmt.Println("[ERROR]: ", err)
		return
	}
}
