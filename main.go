package main

import (
	"log"
	"net/http"
	"strconv"
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
	mux.Handle("/assets/", http.StripPrefix("/assets",
		http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("/healthz", Health)
	mux.HandleFunc("/metrics", apiCfg.Metrics)
	mux.HandleFunc("/reset", apiCfg.Reset)

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

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("Ok"))
}

func (cfg *apiConfig) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	hits := cfg.fileServerHits.Load()
	hit := strconv.Itoa(int(hits))
	w.Write([]byte("Hits: "))
	w.Write([]byte(hit))
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	hits := cfg.fileServerHits.Load()
	hit := strconv.Itoa(int(hits))
	w.Write([]byte(hit))
}
