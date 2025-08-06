package config

type Service struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Endpoint string `yaml:"endpoint"`
}
