package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/wiktorsk8/reverse-proxy/internal/config"
	"github.com/wiktorsk8/reverse-proxy/internal/proxy"
)

func IsFilePath(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	arguments := os.Args

	if len(arguments) < 2 {
		log.Fatal("Config file path is required as first argument")
	}
	configPath := arguments[1]

	if !IsFilePath(configPath) {
		log.Fatal("Invalid config file path")
	}

	authConfig := config.LoadAuthConfig()
	fmt.Println(authConfig.JWTSecret)
	proxyConfig, err := config.LoadProxyConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	newProxyRouter := proxy.NewProxyRouter(proxyConfig, authConfig)

	err = http.ListenAndServe(":8000", newProxyRouter)
	if err != nil {
		panic(err)
	}
}
