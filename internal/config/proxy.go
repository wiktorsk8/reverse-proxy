package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ProxyConfig struct {
	Services  []Service `yaml:"services"`
	RateLimit RateLimit `yaml:"rate-limit"`
}

func LoadProxyConfig(path string) (ProxyConfig, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return ProxyConfig{}, err
	}
	fmt.Println(string(yamlFile))
	conf := &ProxyConfig{}

	err = yaml.Unmarshal(yamlFile, conf)

	if err != nil {
		return ProxyConfig{}, err
	}

	return *conf, nil
}
