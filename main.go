package main

import (
	"log"
	"net/http"

	"github.com/jkvyff/simple-server/internal/database"
)

type apiConfig struct {
	fileserverHits int
	totalRequests  int
	DB *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}  

	cfg := &apiConfig{
		fileserverHits: 0,
		DB: db,
	} 
	
	mux := http.NewServeMux()

	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/*", cfg.middlewareMetrics(fileServer))
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/metrics", cfg.metricsHandler)

	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpId}", cfg.handlerChirpsRetrieveByID)

	mux.HandleFunc("/api/reset", cfg.resetHandler)

	corsMux := middlewareCors(mux)

    server := &http.Server{
        Addr:    ":" + port,
        Handler: corsMux, 
    }

    log.Printf("Starting server on port: %s\n", port)
    log.Fatal(server.ListenAndServe())
}
