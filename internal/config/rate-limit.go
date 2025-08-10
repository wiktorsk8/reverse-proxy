package config

type RateLimit struct {
	Rate  int `yaml:"rate"`
	Burst int `yaml:"burst"`
}
