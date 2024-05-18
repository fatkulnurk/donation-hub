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
	mux.HandleFunc("POST /users/register", api.HandleRegister)
	mux.HandleFunc("POST /users/login", api.HandleLogin)
	mux.HandleFunc("GET /users", api.HandleListUser)
	mux.HandleFunc("GET /projects/upload", api.HandleRequestUploadUrl)
	mux.HandleFunc("POST /projects", api.HandleSubmitProject)
	mux.HandleFunc("PUT /projects/{project_id}/review", api.HandleReviewProjectByAdmin)
	mux.HandleFunc("GET /projects", api.HandleListProject)
	mux.HandleFunc("GET /projects/{project_id}", api.HandleGetProjectById)
	mux.HandleFunc("POST /projects/{project_id}/donations", api.HandleDonateToProject)
	mux.HandleFunc("GET /projects/{project_id}/donations", api.HandleListProjectDonation)
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

func (api *API) HandleRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("register")

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

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("login")
}

func (api *API) HandleListUser(w http.ResponseWriter, r *http.Request) {
	log.Println("get users")
}

func (api *API) HandleRequestUploadUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(NewNotFound())
		return
	}

	log.Println("ini api upload")
}

func (api *API) HandleSubmitProject(w http.ResponseWriter, r *http.Request) {
	log.Println("ini api submit project")
	// Submit Project
	return
}

func (api *API) HandleReviewProjectByAdmin(w http.ResponseWriter, r *http.Request) {
	log.Println("ini api review project by admin")
	// Submit Project
	return
}

func (api *API) HandleListProject(w http.ResponseWriter, r *http.Request) {
	log.Println("ini api list project")
	return
}

func (api *API) HandlePutProjectReview(w http.ResponseWriter, r *http.Request) {

	log.Printf("Review page for project")
}

func (api *API) HandleGetProjectById(w http.ResponseWriter, r *http.Request) {
	log.Printf("Ini project by id")
}

func (api *API) HandleDonateToProject(w http.ResponseWriter, r *http.Request) {
	log.Printf("Ini donate to project")
}

func (api *API) HandleListProjectDonation(w http.ResponseWriter, r *http.Request) {
	log.Printf("Ini handle list project")
}
