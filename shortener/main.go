package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "url required", http.StatusBadRequest)
		return
	}
	code := make([]byte, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	json.NewEncoder(w).Encode(map[string]string{"code": string(code)})
}

func main() {
	http.HandleFunc("POST /shorten", shorten)
	log.Println("shortener on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
