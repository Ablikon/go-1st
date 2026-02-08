package main

import (
	"log"
	"net/http"

	"github.com/Ablikon/go-1st/internal/handlers"
	"github.com/Ablikon/go-1st/internal/middleware"
	"github.com/Ablikon/go-1st/internal/store"
)

func main() {
	st := store.New()
	taskHandler := &handlers.TaskHandler{Store: st}

	mux := http.NewServeMux()
	mux.Handle("/tasks", taskHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	})

	handler := middleware.Logging("request")(
		middleware.APIKey(middleware.DefaultAPIKey)(mux),
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("server started on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
