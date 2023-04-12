package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func (c *Config) FromYaml(path string) (err error) {
	confFile, err := os.Open(path)
	defer confFile.Close()
	data, err := io.ReadAll(confFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, c)
	return
}
