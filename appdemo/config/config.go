package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Conf struct {
	Mysql struct {
		Host     string `yaml:"host"`
		Port     int64  `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}

	Redis struct {
		Host     string `yaml:"host"`
		Port     int64  `yaml:"port"`
		DB       int    `yaml:"db"`
		Password string `yaml:"password"`
	}
}

var V *Conf

func init() {
	f, err := os.Open(`config.yaml`)
	if err != nil {
		panic(err)
	}
	err = yaml.NewDecoder(f).Decode(&V)
	if err != nil {
		panic(err)
	}
}
