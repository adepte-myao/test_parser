package config

type Config struct {
	Server struct {
		BindAddr string `yaml:"bind_addr"`
		LogLevel string `yaml:"log_level"`
	} `yaml:"server"`
}