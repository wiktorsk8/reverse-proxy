package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/wiktorsk8/reverse-proxy/internal/config"
	"github.com/wiktorsk8/reverse-proxy/internal/middleware"
)

func NewProxyRouter(config config.ProxyConfig) *chi.Mux {
	r := chi.NewRouter()

	rateLimiter := middleware.NewRateLimiterMiddleware(config.RateLimit)
	r.Use(rateLimiter.GetMiddleware())

	for _, service := range config.Services {
		handler := getServiceProxyHandler(service)
		r.Handle(service.Endpoint+"/*", handler)
	}

	return r
}

func getServiceProxyHandler(service config.Service) http.Handler {
	reverseProxy := createReverseProxy(service.Host)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, service.Endpoint)
		fmt.Println(r.URL.Path)

		reverseProxy.ServeHTTP(w, r)
	})
}

func createReverseProxy(serviceUrl string) *httputil.ReverseProxy {
	target, err := url.Parse(serviceUrl)
	fmt.Println(target)
	if err != nil {
		panic(err) //TODO: Handle
	}
	return httputil.NewSingleHostReverseProxy(target)
}
