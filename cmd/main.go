package main

import (
	"github.com/wiktorsk8/reverse-proxy/internal/config"
	"github.com/wiktorsk8/reverse-proxy/internal/proxy"
)

func main() {
	proxyConfig := config.LoadProxyConfig()

	proxy := proxy.NewProxy(proxyConfig)

}
