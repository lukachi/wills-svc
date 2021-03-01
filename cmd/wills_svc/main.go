package main

import (
	"github.com/lukachi/wills-svc/internal/app/wills_svc/service"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	vi                     = viper.New()
	config *service.Config = service.NewConfig()
)

func main() {
	cfgPath, exist := os.LookupEnv("CONFIG_FILE")
	if !exist {
		log.Fatal("Config file not found")
	}
	vi.SetConfigName("config")
	vi.SetConfigType("yaml")
	vi.SetConfigFile(cfgPath)
	if err := vi.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := vi.Unmarshal(config); err != nil {
		log.Fatal(err)
	}
	if err := service.Start(config); err != nil {
		log.Fatal(err)
	}
}
