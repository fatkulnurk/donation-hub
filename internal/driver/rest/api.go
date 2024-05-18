package rest

import (
	"encoding/json"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/isdzulqor/donation-hub/internal/core/service/project"
	"github.com/isdzulqor/donation-hub/internal/core/service/user"
	"github.com/isdzulqor/donation-hub/internal/driver/rest/req"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"strconv"
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

	// validation input
	isFails := validator.Validate(rb)
	if isFails != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(isFails.Error()))
		return
	}

	log.Println("ok sekarang register")
	input := model.UserRegisterInput{
		Username: rb.Username,
		Email:    rb.Email,
		Password: rb.Password,
		Role:     rb.Role,
	}

	// store data
	data, err := api.UserService.Register(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(data))
	return
}

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("login")

	var rb req.LoginReqBody
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	// validation input
	isFails := validator.Validate(rb)
	if isFails != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(isFails.Error()))
		return
	}

	data, err := api.UserService.Login(r.Context(), model.UserLoginInput{
		Username: rb.Username,
		Password: rb.Password,
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(NewError(false, "ERR_INVALID_CREDS", "Invalid username or password"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(data))
	return
}

func (api *API) HandleListUser(w http.ResponseWriter, r *http.Request) {
	log.Println("get users")
	var limit, page = 10, 1
	limitQuery := r.URL.Query().Get("limit")
	pageQuery := r.URL.Query().Get("page")

	if limitQuery != "" {
		limit, _ = strconv.Atoi(limitQuery)
	}

	if pageQuery != "" {
		page, _ = strconv.Atoi(pageQuery)
	}

	var role = strings.ToLower(r.URL.Query().Get("role"))
	if role != "donor" && role != "requester" {
		role = ""
	}

	input := model.ListUserInput{
		Limit: int64(limit),
		Page:  int64(page),
		Role:  role,
	}
	data, err := api.UserService.ListUser(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(data))
	return
}

func (api *API) HandleRequestUploadUrl(w http.ResponseWriter, r *http.Request) {
	var mimeType = strings.ToLower(r.URL.Query().Get("mime_type"))
	var fileSize, _ = strconv.Atoi(r.URL.Query().Get("file_size"))

	input := model.RequestUploadUrlInput{
		MimeType: mimeType,
		FileSize: int64(fileSize),
	}
	data, err := api.ProjectService.RequestUploadUrl(r.Context(), input)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(data))
	return
}

func (api *API) HandleSubmitProject(w http.ResponseWriter, r *http.Request) {
	log.Println("ini api submit project")

	var rb req.SubmitProjectReqBody
	err := json.NewDecoder(r.Body).Decode(&rb)

	// error failed parse json
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	// error validation input
	isFails := validator.Validate(rb)
	if isFails != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(isFails.Error()))
		return
	}

	input := model.SubmitProjectInput{
		Title:        rb.Title,
		Description:  rb.Description,
		ImageURLs:    rb.ImageUrls,
		DueAt:        rb.DueAt,
		TargetAmount: rb.TargetAmount,
		Currency:     rb.Currency,
	}
	data, err := api.ProjectService.SubmitProject(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(data))
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
