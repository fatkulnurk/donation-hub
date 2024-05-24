package rest

import (
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/isdzulqor/donation-hub/internal/core/service/auth_token"
	"github.com/isdzulqor/donation-hub/internal/core/service/project"
	"github.com/isdzulqor/donation-hub/internal/core/service/user"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"strings"
)

type API struct {
	DB             *sqlx.DB
	UserService    user.Service
	ProjectService project.Service
	AuthToken      auth_token.Service
}

type customServeMux struct {
	*http.ServeMux
	excludedURLs []string
}

func (api *API) ListenAndServe(cfg *model.ConfigMap) {
	r := &customServeMux{
		ServeMux: http.NewServeMux(),
		// exclude url for force write header Content-Type application/json
		excludedURLs: []string{
			"/",
			"favicon.ico",
			"/assets",
		},
	}
	r.HandleFunc("POST /users/register", api.HandleRegister)
	r.HandleFunc("POST /users/login", api.HandleLogin)
	r.HandleFunc("GET /users", api.HandleListUser)
	r.HandleFunc("GET /projects/{project_id}", api.HandleGetProjectById)
	r.HandleFunc("GET /projects/upload", authTokenMiddleware(api.HandleRequestUploadUrl, api))
	r.HandleFunc("POST /projects", authTokenMiddleware(api.HandleSubmitProject, api))
	r.HandleFunc("GET /projects", authTokenMiddleware(api.HandleListProject, api))
	r.HandleFunc("PUT /projects/{project_id}/review", authTokenMiddleware(api.HandleReviewProjectByAdmin, api))
	r.HandleFunc("POST /projects/{project_id}/donations", authTokenMiddleware(api.HandleDonateToProject, api))
	r.HandleFunc("GET /projects/{project_id}/donations", api.HandleListProjectDonation)
	log.Println("Starting server on :" + cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, r)
	log.Fatal(err)
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
