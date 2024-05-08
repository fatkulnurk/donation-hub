package main

import (
	"encoding/json"
	"github.com/isdzulqor/donation-hub/internal/driver/rest"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	Message string `json:"message"`
}

type customServeMux struct {
	*http.ServeMux
	excludedURLs []string
}

func toJsonMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func (m *customServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestHandler, _ := m.ServeMux.Handler(r)
	found := false
	for _, url := range m.excludedURLs {
		if strings.HasPrefix(r.URL.Path, url) {
			found = true
			break
		}
	}

	// Apply middleware if not excluded and handler exists.
	if !found && requestHandler != nil {
		middleware := toJsonMiddleware(http.HandlerFunc(requestHandler.ServeHTTP))
		middleware(w, r)
	} else if !found {
		http.NotFound(w, r)
	} else {
		requestHandler.ServeHTTP(w, r)
	}
}

func main() {
	mux := &customServeMux{
		ServeMux: http.NewServeMux(),
		excludedURLs: []string{
			"favicon.ico",
		},
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("/users/register", handlePostUserRegister)
	mux.HandleFunc("/users/login", handlePostUserLogin)
	mux.HandleFunc("/users", handleGetUser)
	mux.HandleFunc("/projects/upload", handleGetProjectUpload)
	mux.HandleFunc("/projects/", handleGetAndPostProject)
	mux.HandleFunc("/projects/:id", handleGetProjectById)
	mux.HandleFunc("/projects/:id/review", handlePutProjectReview)
	mux.HandleFunc("/projects/:id/donations", handleGetAndPostProjectDonation)
	log.Println("Starting server on :8180")
	err := http.ListenAndServe(":8180", mux)
	log.Fatal(err)
}

func handlePostUserRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(rest.NewNotFound())
		return
	}
}

func handlePostUserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(rest.NewNotFound())
		return
	}
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(rest.NewNotFound())
		return
	}
}

func handleGetProjectUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(rest.NewNotFound())
		return
	}
}

func handleGetAndPostProject(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// list project
		return
	case http.MethodPost:
		// Submit Project
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(rest.NewNotFound())
		return
	}
}

// PUT: /projects/{project_id}/review
func handlePutProjectReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(rest.NewNotFound())
		return
	}

	projectId := r.URL.Query().Get("id")
	_ = projectId // sementara di taruh di sampah
}

// Get Project
func handleGetProjectById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(rest.NewNotFound())
		return
	}

	projectId := r.URL.Query().Get("id")
	_ = projectId // sementara di taruh di sampah
}

func handleGetAndPostProjectDonation(w http.ResponseWriter, r *http.Request) {
	projectId := r.URL.Query().Get("id")
	_ = projectId // sementara di taruh di sampah

	switch r.Method {
	case http.MethodGet:
		// List Project Donations
		return
	case http.MethodPost:
		// Donate to Project
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(rest.NewNotFound())
		return
	}
}
