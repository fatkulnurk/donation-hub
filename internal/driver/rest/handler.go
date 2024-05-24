package rest

import (
	"encoding/json"
	"fmt"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/isdzulqor/donation-hub/internal/driver/rest/req"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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

	w.Header().Set("Content-Type", "application/json")

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
		UserID:       r.Context().Value("auth_id").(int64),
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

	var rb req.ReviewProjectReqBody
	err := json.NewDecoder(r.Body).Decode(&rb)
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

	projectId, err := strconv.ParseInt(r.PathValue("project_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	// Submit Project
	_, err = api.ProjectService.ReviewProjectByAdmin(r.Context(), model.ReviewProjectByAdminInput{
		UserID:    r.Context().Value("auth_id").(int64),
		ProjectId: projectId,
		Status:    rb.Status,
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(nil))

	return
}

func (api *API) HandleListProject(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if r.URL.Query().Get("limit") != "" {
		limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	}
	startTs, _ := strconv.Atoi(r.URL.Query().Get("start_ts"))
	endTs, _ := strconv.Atoi(r.URL.Query().Get("end_ts"))

	input := model.ListProjectInput{
		UserID:  r.Context().Value("auth_id").(int64),
		Status:  r.URL.Query().Get("status"),
		Limit:   int64(limit),
		StartTs: int64(startTs),
		EndTs:   int64(endTs),
		LastKey: r.URL.Query().Get("last_key"),
	}

	data, err := api.ProjectService.ListProject(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(data))

	return
}

func (api *API) HandleGetProjectById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleGetProjectById")
	var projectId int
	var err error

	// i don't know why the r.PathValue, sometime value is empty
	if r.PathValue("project_id") != "" {
		projectId, err = strconv.Atoi(r.PathValue("project_id"))
		fmt.Println(projectId)
		fmt.Println(err)
	} else {
		projectId, err = strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
		fmt.Println(projectId)
		fmt.Println(err)
	}
	fmt.Println("ProjectID:")
	fmt.Println(projectId)

	input := model.GetProjectByIdInput{
		ProjectId: int64(projectId),
	}

	data, err := api.ProjectService.GetProjectById(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(data))

	return
}

func (api *API) HandleDonateToProject(w http.ResponseWriter, r *http.Request) {
	log.Println("-------------------------------")
	log.Printf("Ini donate to project")
	log.Println("-------------------------------")
	var projectId int
	var err error

	// i don't know why the r.PathValue, sometime value is empty
	if r.PathValue("project_id") != "" {
		projectId, err = strconv.Atoi(r.PathValue("project_id"))
		fmt.Println(projectId)
		fmt.Println(err)
	} else {
		projectId, err = strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
		fmt.Println(projectId)
		fmt.Println(err)
	}

	var rb req.DonateToProjectReqBody
	err = json.NewDecoder(r.Body).Decode(&rb)
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

	input := model.DonateToProjectInput{
		UserID:    r.Context().Value("auth_id").(int64),
		ProjectId: int64(projectId),
		Amount:    rb.Amount,
		Currency:  rb.Currency,
		Message:   rb.Message,
	}

	fmt.Println(input)

	ok, err := api.ProjectService.DonateToProject(r.Context(), input)

	if !ok {
		if err.Error() == "ERR_TOO_MUCH_DONATION" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(NewError(
				false,
				"ERR_TOO_MUCH_DONATION",
				"Donation amount must be less than target amount"))
			return
		} else if err.Error() == "ERR_FORBIDDEN_ACCESS" {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(NewForbiddenAccess())
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest("failed donate to project"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(nil))
	return
}

func (api *API) HandleListProjectDonation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-------------------------------")
	fmt.Println("list project donation")
	fmt.Println("-------------------------------")

	var projectId int
	var err error

	// i don't know why the r.PathValue, sometime value is empty
	if r.PathValue("project_id") != "" {
		projectId, err = strconv.Atoi(r.PathValue("project_id"))
		fmt.Println(projectId)
		fmt.Println(err)
	} else {
		projectId, err = strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
		fmt.Println(projectId)
		fmt.Println(err)
	}
	fmt.Println("ProjectID:")
	limit := 10
	if r.URL.Query().Get("limit") != "" {
		limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	}
	input := model.ListProjectDonationInput{
		ProjectId: int64(projectId),
		Limit:     int64(limit),
		LastKey:   r.URL.Query().Get("last_key"),
	}
	data, err := api.ProjectService.ListDonationByProjectId(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(NewBadRequest(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(NewResponseOk(data))

	return
}
