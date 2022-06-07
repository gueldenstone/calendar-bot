package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Homeserver string   `yaml:"homeserver"`
	Rooms      []string `yaml:"rooms"`
	Calendar   string   `yaml:"calendarURL"`
	NotifyTime string   `yaml:"nofifyTime"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
}

func Parse(path string) (Config, error) {
	conf := Config{}
	yml, err := os.ReadFile(path)
	if err != nil {
		return conf, fmt.Errorf("Reading config file: %w", err)
	}
	err = yaml.Unmarshal(yml, &conf)
	if err != nil {
		return conf, fmt.Errorf("Reading config file: %w", err)
	}
	return conf, nil
}
