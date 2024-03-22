package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
	totalRequests  int 
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	cfg := &apiConfig{} 
	mux := http.NewServeMux()

	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", cfg.middlewareMetrics(fileServer))

	mux.HandleFunc("/healthz", readinessHandler)
	mux.HandleFunc("/metrics", cfg.metricsHandler)
	mux.HandleFunc("/reset", cfg.resetHandler)

	corsMux := middlewareCors(mux)

    server := &http.Server{
        Addr:    ":" + port,
        Handler: corsMux, 
    }

    log.Printf("Starting server on port: %s\n", port)
    log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.totalRequests++
		if r.URL.Path == "/app/" {
			cfg.fileserverHits++
		}

		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprintf(w, "Hits: %d\n", cfg.totalRequests)
}
