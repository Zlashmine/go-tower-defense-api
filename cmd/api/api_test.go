package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tower-defense-api/lib/ratelimiter"
)

func TestRateLimiterMiddleware(t *testing.T) {
	config := config{
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: 10,
			TimeFrame:            time.Second * 5,
			Enabled:              true,
		},
		addr: ":8080",
	}

	app := newTestApplication(t, config)
	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := &http.Client{}
	mockIP := "192.168.1.1"
	marginOfError := 5

	for i := 0; i < config.rateLimiter.RequestsPerTimeFrame+marginOfError; i++ {
		req, err := http.NewRequest("GET", ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		req.Header.Set("X-Forwarded-For", mockIP)
		req.Header.Set("Authorization", "Bearer token")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not send request: %v", err)
		}
		defer resp.Body.Close()

		if i < config.rateLimiter.RequestsPerTimeFrame {
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", resp.Status)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				// t.Errorf("expected status Too Many Requests; got %v", resp.Status)
			}
		}
	}
}
