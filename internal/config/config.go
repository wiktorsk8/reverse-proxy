package config

type ProxyConfig struct {
	Port      string
	Services  []Service
	RateLimit RateLimit
}

func LoadProxyConfig() ProxyConfig {

	return ProxyConfig{}
}
