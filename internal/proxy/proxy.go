package proxy

import (
	"github.com/go-chi/chi/v5"
	"github.com/wiktorsk8/reverse-proxy/internal/config"
	"net/http/httputil"
	"net/url"
)

func NewProxy(config config.ProxyConfig) {
	r := chi.NewRouter()
	setGeneralMiddleware(r)
	loadProxyServices(config.Services, r)
}

func setGeneralMiddleware(r *chi.Mux) {

}

func loadProxyServices(services []config.Service, r *chi.Mux) {
	for _, service := range services {
		reverseProxyInstance := bootServiceReverseProxy(service.Host)
		r.Route(service.Endpoint, func(r chi.Router) {
			r.Mount("/", reverseProxyInstance)
		})
	}
}

func bootServiceReverseProxy(serviceUrl string) *httputil.ReverseProxy {
	target, err := url.Parse(serviceUrl)
	if err != nil {
		panic(err) //TODO: Handle
	}
	return httputil.NewSingleHostReverseProxy(target)
}
