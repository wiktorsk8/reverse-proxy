package tools

import (
	"log"
	"net"
	"net/http"
)

func GetIpFromRequest(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	return host
}
