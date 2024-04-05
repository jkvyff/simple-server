package main

import (
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

	cfg := &apiConfig{
		fileserverHits: 0,
	} 
	
	mux := http.NewServeMux()

	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/*", cfg.middlewareMetrics(fileServer))
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/metrics", cfg.metricsHandler)

	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	mux.HandleFunc("/api/reset", cfg.resetHandler)

	corsMux := middlewareCors(mux)

    server := &http.Server{
        Addr:    ":" + port,
        Handler: corsMux, 
    }

    log.Printf("Starting server on port: %s\n", port)
    log.Fatal(server.ListenAndServe())
}
