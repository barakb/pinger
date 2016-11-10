package pinger

import (
	"time"
)

type Config struct {
	Port      int `json:"port"`
	Addresses []string `json:"addresses"`
	Timeout   time.Duration  `json:"timeout"`
	WaitTime  time.Duration  `json:"wait_time"`
}

func (cfg *Config) AddMachine(name string) {
	for _, s := range cfg.Addresses {
		if s == name {
			return
		}
	}
	cfg.Addresses = append(cfg.Addresses, name)
}

var config *Config

func InitConfig(cfg *Config) {
	config = cfg
	if config.Timeout < 1 {
		config.Timeout = 1
	}
	if config.WaitTime < 1 {
		config.WaitTime = 1
	}
	if config.Port == 0{
		config.Port = 8080
	}
	config.Timeout *= time.Second
	config.WaitTime *= time.Second
}





