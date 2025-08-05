package config

import "os"

type ProxyConfig struct {
	Port      string
	Services  []Service
	RateLimit RateLimit
}

func LoadProxyConfig() ProxyConfig {
	var services []Service
	service := Service{
		Name: "backend",
		Host: "127.0.0.1",
		Port: "8001",
	}

	services = append(services, service)

	return ProxyConfig{
		Port:      os.Getenv("PORT"),
		Services:  services,
		RateLimit: RateLimit{},
	}
}
