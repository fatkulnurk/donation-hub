package rest

import (
	"encoding/json"
	"github.com/isdzulqor/donation-hub/internal/core/service/project"
	"github.com/isdzulqor/donation-hub/internal/core/service/user"
	"github.com/isdzulqor/donation-hub/internal/driver/rest/req"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type API struct {
	DB             *sqlx.DB
	UserService    user.Service
	ProjectService project.Service
}

type customServeMux struct {
	*http.ServeMux
	excludedURLs []string
}

func (api *API) ListenAndServe(appPort string) {
	mux := &customServeMux{
		ServeMux: http.NewServeMux(),
		// exclude url for force write header Content-Type application/json
		excludedURLs: []string{
			"/",
			"favicon.ico",
			"/assets",
		},
	}
	mux.HandleFunc("/users/register", api.HandlePostUserRegister)
	mux.HandleFunc("/users/login", api.HandlePostUserLogin)
	mux.HandleFunc("/users", api.HandleGetUser)
	mux.HandleFunc("/projects/upload", api.HandleGetProjectUpload)
	mux.HandleFunc("/projects", api.HandleGetAndPostProject)
	mux.HandleFunc("/projects/", api.HandleGetProjectById)
	mux.HandleFunc("/projects/{project_id}/review", api.HandlePutProjectReview)
	mux.HandleFunc("/projects/{project_id}/donations", api.HandleGetAndPostProjectDonation)
	log.Println("Starting server on :" + appPort)
	err := http.ListenAndServe(":"+appPort, mux)
	log.Fatal(err)
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

func (api *API) HandlePostUserRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("register")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(NewNotFound())
		return
	}

	var rb req.RegisterReqBody
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	// validation
	isFails := validator.Validate(rb)
	if isFails != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(isFails.Error()))
		return
	}

	log.Println("ok sekarang register")
	// store data
	u, err := api.UserService.Register(rb)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	var data = struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(data))
	return
}

func (api *API) HandlePostUserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(NewNotFound())
		return
	}
	log.Println("login")
}

func (api *API) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(NewNotFound())
		return
	}
	log.Println("get users")
}

func (api *API) HandleGetProjectUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(NewNotFound())
		return
	}

	log.Println("ini api upload")
}

func (api *API) HandleGetAndPostProject(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("ini api list project")
		// list project
		return
	case http.MethodPost:
		log.Println("ini api submit project")
		// Submit Project
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(NewNotFound())
		return
	}
}

// HandlePutProjectReview PUT: /projects/{project_id}/review
func (api *API) HandlePutProjectReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(NewNotFound())
		return
	}

	projectID := extractProjectID(r.URL.Path)
	log.Printf("Review page for project %s\n", projectID)
}

// HandleGetProjectById Get Project
func (api *API) HandleGetProjectById(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodGet {
	//	w.WriteHeader(http.StatusNotFound)
	//	_ = json.NewEncoder(w).Encode(NewNotFound())
	//	return
	//}

	projectId := extractProjectID(r.URL.Path)
	log.Printf("Ini project by id  %s\n", projectId)
}

func (api *API) HandleGetAndPostProjectDonation(w http.ResponseWriter, r *http.Request) {
	projectID := extractProjectID(r.URL.Path)
	log.Printf("Donations page for project %s\n", projectID)
	//
	//switch r.Method {
	//case http.MethodGet:
	//	// List Project Donations
	//	return
	//case http.MethodPost:
	//	// Donate to Project
	//	return
	//default:
	//	w.WriteHeader(http.StatusNotFound)
	//	_ = json.NewEncoder(w).Encode(NewNotFound())
	//	return
	//}
}

func extractProjectID(path string) string {
	re := regexp.MustCompile(`/projects/(\d+)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
