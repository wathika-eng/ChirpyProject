package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	// atomic ensures safety increment across multiple go routines
	fileServerHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	apiCfg := apiConfig{}
	// routes
	mux.Handle("/app/", apiCfg.middlewareMetrics(http.StripPrefix("/app",
		http.FileServer(http.Dir("./templates")))))
	mux.Handle("/app/assets/", apiCfg.middlewareMetrics(http.StripPrefix("/app/assets",
		http.FileServer(http.Dir("./assets")))))

	mux.HandleFunc("GET /api/healthz", Health)
	mux.HandleFunc("GET /admin/metrics", apiCfg.Metrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.Reset)
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirp)

	log.Printf("server running on http://localhost%v\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("error running server")
	}
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(+1)
		next.ServeHTTP(w, r)
	}
}
