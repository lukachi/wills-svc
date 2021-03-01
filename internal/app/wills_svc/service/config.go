package service

type Config struct {
	LogLevel     string `mapstructure:"log_level"`
	DatabaseURL  string `mapstructure:"db_url"`
	ListenerAddr string `mapstructure:"listener_addr"`
}

func NewConfig() *Config {
	return &Config{}
}
