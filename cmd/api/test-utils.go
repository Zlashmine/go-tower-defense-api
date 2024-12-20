package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"tower-defense-api/lib/ratelimiter"
	"tower-defense-api/lib/repository"
	"tower-defense-api/lib/repository/cache"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	cfg.authToken = "token"

	logger := zap.Must(zap.NewDevelopment()).Sugar()
	mockRepository := repository.NewMockRepository()
	mockCacheStore := cache.NewMockStore()

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	return &application{
		logger:     logger,
		repository: mockRepository,
		cacheStore: mockCacheStore,
		rateLimiter: rateLimiter,
		config:     cfg,
	}
}

func executeRequest(mux http.Handler, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer token")
	mux.ServeHTTP(recorder, req)
	return recorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected status code %d, got %d", expected, actual)
	}
}
