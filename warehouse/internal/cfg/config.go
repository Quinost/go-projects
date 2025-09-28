package cfg

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerCfg   `yaml:"server"`
	Database DatabaseCfg `yaml:"database"`
	Auth     AuthCfg     `yaml:"auth"`
}

type DatabaseCfg struct {
	Driver           string `yaml:"driver"`
	ConnectionString string `yaml:"connection_string"`
	Seed             bool   `yaml:"seed"`
}

type ServerCfg struct {
	Port string `yaml:"port"`
}

type AuthCfg struct {
	Secret string `yaml:"secret"`
	Exp    int    `yaml:"exp_min"`
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
