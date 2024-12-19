package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"tower-defense-api/lib/ratelimiter"
	"tower-defense-api/lib/repository/cache"
)

func TestGetUser(t *testing.T) {
	config := config{
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: 20,
			TimeFrame:            time.Second * 5,
			Enabled:              false, // ,true
		},
		addr: ":8080",
	}

	app := newTestApplication(t, config)
	mux := app.mount()

	t.Run("should return a user by Id", func(t *testing.T) {
		mockCacheStore := app.cacheStore.Users.(*cache.MockUsersStore)

		mockCacheStore.On("Get", int64(1)).Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		recorder := executeRequest(mux, req)
		checkResponseCode(t, http.StatusOK, recorder.Code)

		mockCacheStore.AssertNumberOfCalls(t, "Get", 0) // 1
		mockCacheStore.AssertNumberOfCalls(t, "Set", 0) // 1

		mockCacheStore.Calls = nil
	})
}
