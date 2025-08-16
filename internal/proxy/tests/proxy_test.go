package proxy

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wiktorsk8/reverse-proxy/internal/config"
	proxy2 "github.com/wiktorsk8/reverse-proxy/internal/proxy"
)

func getHostAndPort(targetUrl string) (string, string) {
	parsed, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatal(err)
	}
	host, port, err := net.SplitHostPort(parsed.Host)

	if err != nil {
		log.Fatal(err)
	}

	return host, port
}

func getAuthConfig() config.AuthConfig {
	return config.AuthConfig{
		JWTSecret: "top_secret",
	}
}

func getConfig(services []config.Service) config.ProxyConfig {
	return config.ProxyConfig{
		Services: services,
		RateLimit: config.RateLimit{
			Rate:  1,
			Burst: 5,
		},
	}
}

func getService(host, port, endpoint string) config.Service {
	return config.Service{
		Host:     host,
		Endpoint: endpoint,
		Port:     port,
	}
}

func getServiceMockHandler(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"service":  name,
			"response": "Success",
		})
	})
}

func getBearerToken(authConfig config.AuthConfig) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "user@company.com",
	})

	signedToken, err := token.SignedString([]byte(authConfig.JWTSecret))
	if err != nil {
		log.Fatal(err)
	}

	return "Bearer " + signedToken
}

func TestHelloWorld(t *testing.T) {
	helloWorldService := httptest.NewServer(getServiceMockHandler("helloworld"))
	defer helloWorldService.Close()
	host, port := getHostAndPort(helloWorldService.URL)
	service := getService("http://"+host, port, "/hello-world")

	proxyConfig := getConfig([]config.Service{service})
	authConfig := getAuthConfig()

	proxyRouter := proxy2.NewProxyRouter(proxyConfig, authConfig)
	server := httptest.NewServer(proxyRouter)
	defer server.Close()

	tests := []struct {
		httpMethod string
		service    string
		url        string
	}{
		{"GET", service.Name, server.URL + service.Endpoint + "/smth"},
		{"POST", service.Name, server.URL + service.Endpoint + "/smth-else"},
	}

	token := getBearerToken(authConfig)

	for _, test := range tests {
		req, _ := http.NewRequest(test.httpMethod, test.url, nil)
		req.Header.Set("Authorization", token)

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("Request failed: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			responseBody, _ := io.ReadAll(response.Body)
			t.Errorf("Request failed with status %v %s", response.StatusCode, responseBody)
		}

		responseBody, _ := io.ReadAll(response.Body)

		var resData map[string]string
		json.Unmarshal(responseBody, &resData)

		if resData["service"] != test.service && resData["response"] != "Success" {
			t.Errorf("Request failed for service %s and response %s", test.service, resData["response"])
		}

		response.Body.Close()
	}
}
