package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wiktorsk8/reverse-proxy/internal/config"
	"github.com/wiktorsk8/reverse-proxy/internal/proxy"
)

func IsFilePath(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func main() {
	arguments := os.Args

	if len(arguments) < 2 {
		log.Fatal("Config file path is required as first argument")
	}
	configPath := arguments[1]

	if !IsFilePath(configPath) {
		log.Fatal("Invalid config file path")
	}

	proxyConfig, err := config.LoadProxyConfig(configPath)

	if err != nil {
		log.Fatal(err)
	}

	newProxyRouter := proxy.NewProxyRouter(proxyConfig)

	err = http.ListenAndServe(":8000", newProxyRouter)
	if err != nil {
		panic(err)
	}
}
