package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Homeserver string   `yaml:"homeserver"`
	Calendar   string   `yaml:"calendarURL"`
	NotifyTime string   `yaml:"nofifyTime"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
	Rooms      []string `yaml:"rooms"`
}

func Parse(path string) (Config, error) {
	conf := Config{}
	yml, err := os.ReadFile(path)
	if err != nil {
		return conf, fmt.Errorf("reading config file: %w", err)
	}
	err = yaml.Unmarshal(yml, &conf)
	if err != nil {
		return conf, fmt.Errorf("reading config file: %w", err)
	}
	return conf, nil
}
