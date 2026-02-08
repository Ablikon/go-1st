package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(message string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timestamp := time.Now().Format(time.RFC3339)
			log.Printf("%s %s %s %s", timestamp, r.Method, r.URL.Path, message)
			next.ServeHTTP(w, r)
		})
	}
}
