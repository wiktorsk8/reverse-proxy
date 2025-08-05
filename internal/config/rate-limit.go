package config

type RateLimit struct {
	Timeout         int
	AllowedRequests int
	WindowDuration  int
}
