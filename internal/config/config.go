package config

type Config struct {
	Server ServerConfig `yaml:"server"`
	Store  StoreConfig  `yaml:"store"`
}

type ServerConfig struct {
	BindAddr string `yaml:"bind_addr"`
	LogLevel string `yaml:"log_level"`
	BaseLink string `yaml:"link"`
}

type StoreConfig struct {
	DatabaseURL string `yaml:"database_url"`
}
