package utils

import "net/http"

type ApiConfig struct {
	FileserverHits int
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	cfg.FileserverHits++
	return next
}
