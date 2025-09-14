package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Pages    []Site   `yaml:"sites"`
	Settings Settings `yaml:"settings"`
}

type Site struct {
	URL      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
}

type Settings struct {
	RequestTimeout time.Duration `yaml:"request_timeout"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return &cfg, err
}
