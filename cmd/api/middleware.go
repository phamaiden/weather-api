package main

import (
	"encoding/json"
	"net/http"

	"github.com/didip/tollbooth/v8"
)

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// var lmt = tollbooth.NewLimiter(1, nil)
// lmt.SetBurst(5)

// var lmt = func() *limiter.Limiter {
// 	l := tollbooth.NewLimiter(1, nil)
// 	l.SetBurst(5)
// 	return l
// }()

func limiterMiddleware(next http.Handler) http.Handler {
	lmt := tollbooth.NewLimiter(1, nil)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tollbooth.LimitByRequest(lmt, w, r)
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error: ": err.Error(),
			})
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
