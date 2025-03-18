package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	hits := cfg.fileServerHits.Load()
	w.Write([]byte(fmt.Sprintf(`
			<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	`, hits)))
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	hits := cfg.fileServerHits.Load()
	hit := strconv.Itoa(int(hits))
	w.Write([]byte(hit))
}

// decode
func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type BodyStruct struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	bodyStruct := BodyStruct{}
	err := decoder.Decode(&bodyStruct)
	if err != nil {
		w.Write([]byte("error:"))
		w.Write([]byte("something went wrong"))
		return
	}
	w.Write([]byte("valid:"))
	w.Write([]byte("true"))
}
