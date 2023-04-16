package config

import (
	"fmt"
	"go-backend/pkg/logger"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Server  struct {
		Type   string `yaml:"type"`
		BindIp string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"server"`
}

var singleToneInstance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logger.GetLogger()
		logger.Info("Read App config")
		singleToneInstance := &Config{}
		if err := cleanenv.ReadConfig("config.yml", singleToneInstance); err != nil {
			help, _ := cleanenv.GetDescription(singleToneInstance, nil)
			logger.Info(help)

			logger.Fatal(err)
		}
	})
	fmt.Println(singleToneInstance)
	return singleToneInstance
}
