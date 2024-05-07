package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			resp := Response{"Method salah"}
			_ = json.NewEncoder(w).Encode(resp)

			return
		}

		w.WriteHeader(http.StatusOK)
		resp := Response{"oke handler bisa"}
		_ = json.NewEncoder(w).Encode(resp)
	})
	log.Println("Starting server on :8180")
	err := http.ListenAndServe(":8180", mux)
	log.Fatal(err)
}
